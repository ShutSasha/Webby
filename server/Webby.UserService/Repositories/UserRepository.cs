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

   public async Task<List<UserFollower>> GetUserFollows(Guid userId)
   {
      return await _context.UserFollowers
         .Include(uf => uf.FollowedUser)
         .ThenInclude(uf => uf.Followers)
         .Where(uf => uf.FollowerId == userId)
         .ToListAsync();
   }

   public async Task<List<UserFollower>> GetUserFollowers(Guid userId)
   {
      return await _context.UserFollowers
         .Include(uf => uf.FollowerUser)
         .ThenInclude(u => u.Followers)
         .Where(uf => uf.UserId == userId)
         .ToListAsync();
   }
   
   public async Task<List<User>> GetByIds(List<Guid> ids)
   {
      return await _context.Users
         .Where(u => ids.Contains(u.UserId))
         .ToListAsync();
   }

   public async Task<List<Guid>> GetSubscriptionIds(Guid requestUserId)
   {
      return await _context.UserFollowers
         .Where(uf => uf.FollowerId == requestUserId)
         .Select(uf => uf.UserId)
         .ToListAsync();
   }
}