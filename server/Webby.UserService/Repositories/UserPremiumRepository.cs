using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class UserPremiumRepository : GenericRepository<UserPremium>, IUserPremiumRepository
{
   public UserPremiumRepository(AppDbContext context) : base(context)
   {
   }
}