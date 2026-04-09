using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class UserPremiumRepository : GenericRepository<UserPremium>, IUserPremiumRepository
{
   public UserPremiumRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<UserPremium?> GetUserPremiumInformation(Guid userId) =>
      await _context.UserPremiums
         .Where(up => up.UserId == userId)
         .Include(up => up.User)
         .FirstOrDefaultAsync();
}