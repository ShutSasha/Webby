using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;
using Webby.UserService.Models.Enums;

namespace Webby.UserService.Repositories;

public class PaymentRepository : GenericRepository<Payment>, IPaymentRepository
{
   public PaymentRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<Payment> GetByExternalId(string externalId)
   {
      return await _context.Payments
         .Where(p => p.ExternalId == externalId)
         .FirstAsync();
   }
   
   public async Task<Dictionary<int, decimal>> GetMonthlySubscriptionRevenueAsync(DateTime startDate, DateTime endDate)
   {
      var groupedData = await _context.Payments
         .Where(p => p.CreatedAt >= startDate 
                     && p.CreatedAt < endDate 
                     && p.Status == PaymentStatus.Succeeded)
         .GroupBy(p => new { p.CreatedAt.Year, p.CreatedAt.Month })
         .Select(g => new
         {
            Month = g.Key.Month,
            TotalAmount = g.Sum(p => p.Amount)
         })
         .ToListAsync();

      return groupedData.ToDictionary(x => x.Month, x => x.TotalAmount);
   }

   public async Task<int> GetMonthRevenue()
   {
      return (int)await _context.Payments
         .Where(p => p.CreatedAt >= DateTime.UtcNow)
         .SumAsync(p => p.Amount);
   }
   
}