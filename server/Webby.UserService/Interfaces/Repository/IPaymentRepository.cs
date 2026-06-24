using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IPaymentRepository : IRepository<Payment>
{
   Task<Payment> GetByExternalId(string externalId);
   Task<Dictionary<int, decimal>> GetMonthlySubscriptionRevenueAsync(DateTime startDate, DateTime endDate);
   Task<int> GetMonthRevenue();
   Task<List<Payment>> GetPendingPaymentsOlderThan(DateTime thresholdTime);
   
}