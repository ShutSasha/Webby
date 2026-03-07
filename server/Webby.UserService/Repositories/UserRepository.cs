using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Dtos.User;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class UserRepository : GenericRepository<User>,IUserRepository
{
   public UserRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<UserFollowStats> GetUserFollowBlock(Guid userId)
   {
      return new UserFollowStats
      {
         Followers = await _context.UserFollowers
            .CountAsync(x => x.UserId == userId),

         Following = await _context.UserFollowers
            .CountAsync(x => x.FollowerId == userId)
      };
   }
}