using System.Text.Json;
using StackExchange.Redis;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Event;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services.Handlers;

public class UserDeletedEventHandler : BackgroundService
{
    private readonly IDatabase _redisDatabase;
    private readonly ILogger<UserDeletedEventHandler> _logger;
    private readonly IServiceScopeFactory _scopeFactory;
    private bool _isGroupCreated;
    
    public UserDeletedEventHandler(
        IConnectionMultiplexer redis, 
        ILogger<UserDeletedEventHandler> logger,
        IServiceScopeFactory scopeFactory)
    {
        _redisDatabase = redis.GetDatabase();
        _logger = logger;
        _scopeFactory = scopeFactory;
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
                            using var scope = _scopeFactory.CreateScope();
                            
                            var videoService = scope.ServiceProvider.GetRequiredService<IVideoService>();
                            var playlistService = scope.ServiceProvider.GetRequiredService<IPlaylistService>();

                            await videoService.ClearUserVideos(data.UserId);
                            await playlistService.ClearPlaylists(data.UserId);
                        }
                    }
                    
                    await _redisDatabase.StreamAcknowledgeAsync(RedisConstants.UserStreamName, "video-service-group", message.Id);
                }
                else
                {
                    await Task.Delay(1000, stoppingToken);
                }
            }
            catch (Exception)
            {
                _logger.LogWarning("Failed to connect to Redis or process stream. Retrying in {Timeout}s", RedisConstants.RedisConnectionTimeout);
                await Task.Delay(TimeSpan.FromSeconds(RedisConstants.RedisConnectionTimeout), stoppingToken);
            }
        }
    }
}