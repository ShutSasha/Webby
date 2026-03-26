using Webby.VideoService.Interfaces.Helpers;

namespace Webby.VideoService.Services.Background;

public class QueuedHostedService : BackgroundService
{
   private readonly IBackgroundTaskQueue _taskQueue;
   private readonly ILogger<QueuedHostedService> _logger;

   public QueuedHostedService(
      IBackgroundTaskQueue taskQueue,
      ILogger<QueuedHostedService> logger)
   {
      _taskQueue = taskQueue;
      _logger = logger;
   }

   protected override async Task ExecuteAsync(CancellationToken stoppingToken)
   {
      _logger.LogInformation("Background queue started");

      while (!stoppingToken.IsCancellationRequested)
      {
         try
         {
            var workItem = await _taskQueue.DequeueAsync(stoppingToken);

            await workItem(stoppingToken);
         }
         catch (Exception ex)
         {
            _logger.LogError(ex, "Error occurred executing background task");
         }
      }
   }
}