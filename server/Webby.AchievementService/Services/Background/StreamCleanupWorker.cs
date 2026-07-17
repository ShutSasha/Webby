using StackExchange.Redis;
using Webby.AchievementService.Constants;

namespace Webby.AchievementService.Services.Background;

public class StreamCleanupWorker : BackgroundService
{
   private readonly IConnectionMultiplexer _redis;
   private readonly ILogger<StreamCleanupWorker> _logger;

   public StreamCleanupWorker(IConnectionMultiplexer redis, ILogger<StreamCleanupWorker> logger)
   {
      _redis = redis;
      _logger = logger;
   }

   protected override async Task ExecuteAsync(CancellationToken stoppingToken)
   {
      while (!stoppingToken.IsCancellationRequested)
      {
         try
         {
            var db = _redis.GetDatabase();
            using var timer = new PeriodicTimer(TimeSpan.FromHours(1));

            while (await timer.WaitForNextTickAsync(stoppingToken))
            {
               _logger.LogInformation("Cleaning events in {StreamName}...", RedisConstants.SteamName);
                    
               var trimmedCount = await db.StreamTrimAsync(
                  key: RedisConstants.SteamName, 
                  maxLength: RedisConstants.MaxStreamLength, 
                  useApproximateMaxLength: true);

               if (trimmedCount > 0)
               {
                  _logger.LogInformation("Stream clean is finished");
               }
            }
         }
         catch (RedisConnectionException)
         {
            _logger.LogWarning("Redis connection lost. Timeout {Timeout} ms",RedisConstants.RedisConnectionTimeout);
            await Task.Delay(RedisConstants.RedisConnectionTimeout, stoppingToken);
         }
         catch (OperationCanceledException)
         {
            break;
         }
         catch (Exception ex)
         {
            _logger.LogError(ex, "Redis Stream error. Timeout {Timeout}",RedisConstants.RedisConnectionTimeout);
            await Task.Delay(RedisConstants.RedisConnectionTimeout, stoppingToken);
         }
      }
   }
}