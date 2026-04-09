using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IUserPremiumRepository : IRepository<UserPremium>
{
   Task<UserPremium?> GetUserPremiumInformation(Guid userId);
}