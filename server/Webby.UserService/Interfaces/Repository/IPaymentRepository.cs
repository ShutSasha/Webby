using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IPaymentRepository : IRepository<Payment>
{
   Task<Payment> GetByExternalId(string externalId);
}