using Microsoft.AspNetCore.Mvc;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Service;

public interface IAchievementService
{
   Task<Achievement> CreateAchievement(CreateAchievementRequest request);
   Task DeleteAchievement(Guid achievementId);
   Task<Achievement> UpdateAchievement(UpdateAchievementRequest request);
   Task<List<Achievement>> GetAchievements();
   Task<Achievement?> FindById(Guid achievementId);
   Task AddUserAchievement(Guid userId, Guid achievementId);
   Task<UserAchievement> GetUserAchievement(Guid userId, Guid achievementId);
   Task UpdateUserAchievement(UserAchievement userAchievement);
   Task<List<ProfileAchievementDto>> GetPinnedAchievements(Guid userId);
   Task<GetUserAchievementsResponse> GetUserAchievementsBlock(Guid userId);
   
}