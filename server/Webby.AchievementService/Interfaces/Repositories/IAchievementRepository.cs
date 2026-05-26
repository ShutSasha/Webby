using Webby.AchievementService.Dtos.Achievement;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Interfaces.Repositories;

public interface IAchievementRepository : IRepository<Achievement>
{
   Task AddUserAchievement(Guid userId, Guid achievementId);
   Task<bool> HasUserAchievement(Guid userId, Guid achievementId);
   Task<UserAchievement?> GetUserAchievement(Guid userId, Guid achievementId);
   Task UpdateUserAchievement(UserAchievement userAchievement);
   Task<List<UserAchievement>> GetUserPinnedAchievements(Guid userId);
   Task<List<AchievementDto>> GetAchievementsWithUserStatus(Guid userId);
   Task<int> CountPinnedAchievements(Guid userId);
}