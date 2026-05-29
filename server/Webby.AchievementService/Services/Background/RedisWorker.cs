using StackExchange.Redis;
using Webby.AchievementService.Constants;
using Webby.AchievementService.Services.Handlers;

namespace Webby.AchievementService.Services.Background;

public class RedisWorker : BackgroundService
{
    private readonly IConnectionMultiplexer _redis;
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly ILogger<RedisWorker> _logger;

    private const string StreamName = "events:platform";
    private const string GroupName = "achievement-service-group";
    private const string ConsumerName = "worker-1";

    public RedisWorker(IConnectionMultiplexer redis, IServiceScopeFactory scopeFactory, ILogger<RedisWorker> logger)
    {
        _redis = redis;
        _scopeFactory = scopeFactory;
        _logger = logger;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            try
            {
                var db = _redis.GetDatabase();

                try
                {
                    await db.StreamCreateConsumerGroupAsync(StreamName, GroupName, StreamPosition.NewMessages);
                }
                catch (RedisServerException ex) when (ex.Message.Contains("BUSYGROUP"))
                {
                }

                while (!stoppingToken.IsCancellationRequested)
                {
                    var messages = await db.StreamReadGroupAsync(StreamName, GroupName, ConsumerName, ">", count: 10);

                    foreach (var msg in messages)
                    {
                        var type = msg.Values.FirstOrDefault(v => v.Name == "type").Value.ToString();
                        var payload = msg.Values.FirstOrDefault(v => v.Name == "payload").Value.ToString();

                        using (var scope = _scopeFactory.CreateScope())
                        {
                            var handler = scope.ServiceProvider.GetRequiredService<AchievementHandler>();
                            await handler.HandleEventAsync(type, payload);
                        }

                        await db.StreamAcknowledgeAsync(StreamName, GroupName, msg.Id);
                    }

                    if (messages.Length == 0)
                    {
                        await Task.Delay(500, stoppingToken);
                    }
                }
            }
            catch (RedisConnectionException ex)
            {
                _logger.LogWarning("Redis connection lost. Timeout {Timeout} ms", RedisConstants.RedisConnectionTimeout);
                await Task.Delay(RedisConstants.RedisConnectionTimeout, stoppingToken);
            }
            catch (OperationCanceledException)
            {
                break;
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Worker error. Timeout {Timeout} ms", RedisConstants.RedisConnectionTimeout);
                await Task.Delay(RedisConstants.RedisConnectionTimeout, stoppingToken);
            }
        }
    }
}