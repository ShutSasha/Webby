using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class AchievementRepository : GenericRepository<Achievement>, IAchievementRepository
{
   public AchievementRepository(AppDbContext context) : base(context)
   {
   }

   public async Task AddUserAchievement(Guid userId, Guid achievementId)
   {
      var entity = new UserAchievement()
      {
         AchievementId = achievementId,
         IsPinned = false,
         UnlockedAt = DateTime.UtcNow,
         UserId = userId
      };

      await _context.UserAchievements.AddAsync(entity);
      await _context.SaveChangesAsync();
   }

   public async Task<bool> HasUserAchievement(Guid userId, Guid achievementId)
   {
      return await _context.UserAchievements
         .AnyAsync(ua => ua.UserId == userId && ua.AchievementId == achievementId);
   }

   public async Task<UserAchievement?> GetUserAchievement(Guid userId, Guid achievementId)
   {
      return await _context.UserAchievements
         .Where(ua => ua != null && ua.UserId == userId && ua.AchievementId == achievementId)
         .FirstOrDefaultAsync();
   }

   public async Task UpdateUserAchievement(UserAchievement userAchievement)
   {
      _context.UserAchievements.Update(userAchievement);
      await _context.SaveChangesAsync();
   }
   
   public async Task<List<UserAchievement>> GetUserPinnedAchievements(Guid userId)
   {
      return await _context.UserAchievements
         .Include(ua => ua.Achievement)
         .Where(ua => ua.UserId == userId && ua.IsPinned)
         .Take(3)
         .ToListAsync();
   }
   
   public async Task<List<AchievementDto>> GetAchievementsWithUserStatus(Guid userId)
   {
      return await _context.Achievements
         .Where(a => !_context.UserAchievements
            .Any(ua => ua.UserId == userId 
                       && ua.AchievementId == a.AchievementId 
                       && ua.IsPinned))
         .Select(a => new AchievementDto
         {
            AchievementId = a.AchievementId,
            IconUrl = a.IconUrl,
            Title = a.Title,
            Description = a.Description,
            IsUnlocked = _context.UserAchievements
               .Any(ua => ua.UserId == userId && ua.AchievementId == a.AchievementId)
         })
         .ToListAsync();
   }

   public async Task<int> CountPinnedAchievements(Guid userId) 
      => await _context.UserAchievements
         .CountAsync(a => a.UserId == userId && a.IsPinned == true);
}