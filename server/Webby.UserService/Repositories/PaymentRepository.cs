using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

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
}