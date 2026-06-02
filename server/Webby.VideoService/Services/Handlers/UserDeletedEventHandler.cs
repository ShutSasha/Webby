using System.Text.Json;
using Microsoft.EntityFrameworkCore;
using StackExchange.Redis;
using Webby.VideoService.Constants;
using Webby.VideoService.Data;
using Webby.VideoService.Dtos.Event;

namespace Webby.VideoService.Services.Handlers;

public class UserDeletedEventHandler : BackgroundService
{
    private readonly IDatabase _redisDatabase;
    private readonly IServiceProvider _serviceProvider;
    private readonly ILogger<UserDeletedEventHandler> _logger;
    private bool _isGroupCreated;

    public UserDeletedEventHandler(
        IConnectionMultiplexer redis, 
        IServiceProvider serviceProvider,
        ILogger<UserDeletedEventHandler> logger)
    {
        _redisDatabase = redis.GetDatabase();
        _serviceProvider = serviceProvider;
        _logger = logger;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            try
            {
                if (!_isGroupCreated)
                {
                    try
                    {
                        await _redisDatabase.StreamCreateConsumerGroupAsync(RedisConstants.UserStreamName, "video-service-group", "0-0", true);
                    }
                    catch (RedisServerException ex) when (ex.Message.Contains("BUSYGROUP"))
                    {
                    }
                    
                    _isGroupCreated = true;
                }

                var result = await _redisDatabase.StreamReadGroupAsync(
                    RedisConstants.UserStreamName, 
                    "video-service-group", 
                    "video-worker-1", 
                    ">", 
                    1);
                
                if (result.Length > 0)
                {
                    var message = result.First();
                    var eventType = message.Values.FirstOrDefault(x => x.Name == "EventType").Value;

                    if (eventType == "user_deleted")
                    {
                        var payload = message.Values.FirstOrDefault(x => x.Name == "Payload").Value;
                        var data = JsonSerializer.Deserialize<UserDeletedEventDto>(payload);

                        if (data != null)
                        {
                            using var scope = _serviceProvider.CreateScope();
                            var dbContext = scope.ServiceProvider.GetRequiredService<AppDbContext>();

                            await dbContext.Videos.Where(v => v.UserId == data.UserId).ExecuteDeleteAsync(stoppingToken);
                            await dbContext.Playlists.Where(p => p.UserId == data.UserId).ExecuteDeleteAsync(stoppingToken);
                        }
                    }

                    await _redisDatabase.StreamAcknowledgeAsync(RedisConstants.StreamName, "video-service-group", message.Id);
                }
                else
                {
                    await Task.Delay(1000, stoppingToken);
                }
            }
            catch (Exception)
            {
                _logger.LogWarning( "Failed to connect to Redis or process stream. Retrying in {Timeout}s",RedisConstants.RedisConnectionTimeout);
                await Task.Delay(TimeSpan.FromSeconds(RedisConstants.RedisConnectionTimeout), stoppingToken);
            }
        }
    }
}