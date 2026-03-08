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
   
   public async Task AddUserFollowing(Guid userId, Guid followerId)
   {
      var entity = new UserFollower
      {
         UserId = userId,
         FollowerId = followerId
      };

      await _context.UserFollowers.AddAsync(entity);
      await _context.SaveChangesAsync();
   }

   public async Task DeleteUserFollowing(Guid userId, Guid followerId)
   {
      var userFollowing = await _context.UserFollowers
         .Where(uf => uf.FollowerId == followerId)
         .FirstOrDefaultAsync();

      if (userFollowing != null) _context.UserFollowers.Remove(userFollowing);

      await _context.SaveChangesAsync();
   }

   public async Task<bool> HasUserFollow(Guid userId, Guid followerId)
   {
      return await _context.UserFollowers
         .AnyAsync(uf => uf.FollowerId == followerId && uf.UserId == userId);
   }
}