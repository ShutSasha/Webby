using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IAchievementRepository : IRepository<Achievement>
{
   Task AddUserAchievement(Guid userId, Guid achievementId);
   Task<bool> HasUserAchievement(Guid userId, Guid achievementId);
   Task<UserAchievement?> GetUserAchievement(Guid userId, Guid achievementId);
   Task UpdateUserAchievement(UserAchievement userAchievement);
   Task<List<UserAchievement>> GetUserPinnedAchievements(Guid userId);
   Task<List<AchievementDto>> GetAchievementsWithUserStatus(Guid userId);
}