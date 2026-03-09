using Microsoft.AspNetCore.Mvc;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;

namespace Webby.UserService.Controllers;

[ApiController]
[Route("api/achievements")]
public class AchievementController: ControllerBase
{
   private readonly IAchievementService _achievementService;
   
   public AchievementController(IAchievementService achievementService)
   {
      _achievementService = achievementService;
   }

   [HttpGet]
   public async Task<ActionResult<ApiResponse<List<Achievement>>>> GetAchievements()
   {
      var achievements = await _achievementService.GetAchievements();
      return Ok(ApiResponse<List<Achievement>>.Ok("Successfully extracted achievements", achievements));
   }

   [HttpPost]
   public async Task<ActionResult<ApiResponse<Achievement>>> CreateAchievement([FromForm] CreateAchievementRequest request)
   {
      var achievementCreationResult = await _achievementService.CreateAchievement(request);
      return Ok(ApiResponse<Achievement>.Ok("Successully create achievement", achievementCreationResult));
   }

   [HttpDelete("{achievementId:guid}")]
   public async Task<ActionResult<ApiResponse>> DeleteAchievement([FromRoute] Guid achievementId)
   {
      await _achievementService.DeleteAchievement(achievementId);
      return Ok(ApiResponse.Ok("Successfully deleted achievement"));
   }

   [HttpPut]
   public async Task<ActionResult<ApiResponse<Achievement>>> UpdateAchievement([FromForm] UpdateAchievementRequest request)
   {
      var updateAchievementResult = await _achievementService.UpdateAchievement(request);
      return Ok(ApiResponse<Achievement>.Ok("Successfully updated achievement", updateAchievementResult));
   }
   
}