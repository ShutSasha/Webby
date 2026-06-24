using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Services.Background;

public class PaymentBackgroundWorker : BackgroundService
{
   private readonly IServiceProvider _serviceProvider;
   private readonly TimeSpan _checkInterval = TimeSpan.FromMinutes(15); 

   public PaymentBackgroundWorker(IServiceProvider serviceProvider)
   {
      _serviceProvider = serviceProvider;
   }

   protected override async Task ExecuteAsync(CancellationToken stoppingToken)
   {
      while (!stoppingToken.IsCancellationRequested)
      {
         try
         {
            await TrackAndSyncPayments(stoppingToken);
         }
         catch (Exception ex)
         {
            Console.WriteLine(ex.Message);
         }

         await Task.Delay(_checkInterval, stoppingToken);
      }
   }

   private async Task TrackAndSyncPayments(CancellationToken stoppingToken)
   {
      using var scope = _serviceProvider.CreateScope();
        
      var paymentRepository = scope.ServiceProvider.GetRequiredService<IPaymentRepository>();
      var paymentService = scope.ServiceProvider.GetRequiredService<IPaymentService>();

      var thresholdTime = DateTime.UtcNow.AddMinutes(-15);
      var pendingPayments = await paymentRepository.GetPendingPaymentsOlderThan(thresholdTime);

      if (!pendingPayments.Any()) return;

      foreach (var payment in pendingPayments)
      {
         if (stoppingToken.IsCancellationRequested) break;
         await paymentService.SyncPendingPayment(payment);
      }
   }
}