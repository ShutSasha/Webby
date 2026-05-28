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
      using var timer = new PeriodicTimer(TimeSpan.FromHours(1));
      var db = _redis.GetDatabase();

      try
      {
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
      catch (OperationCanceledException)
      {
         _logger.LogInformation("Redis stream worker is stopped.");
      }
      catch (Exception ex)
      {
         _logger.LogError(ex, "Redis Stream error");
      }
   }
}