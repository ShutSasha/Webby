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
   
}