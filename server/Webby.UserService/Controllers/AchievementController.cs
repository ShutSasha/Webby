using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
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
   [SwaggerOperation("Get all platform achievements")]
   public async Task<ActionResult<ApiResponse<List<Achievement>>>> GetAchievements()
   {
      var achievements = await _achievementService.GetAchievements();
      return Ok(ApiResponse<List<Achievement>>.Ok("Successfully extracted achievements", achievements));
   }

   [HttpGet("{userId:guid}")]
   [SwaggerOperation("Retrieve a user's pinned and unlocked achievements")]
   public async Task<ActionResult<ApiResponse<GetUserAchievementsResponse>>> GetUserAchievementsBlock([FromRoute] Guid userId)
   {
      var userAchievementsWithStatus = await _achievementService.GetUserAchievementsBlock(userId);
      return Ok(ApiResponse<GetUserAchievementsResponse>.Ok("Successfully retrieved", userAchievementsWithStatus));
   }
   
   [HttpPost]
   [SwaggerOperation("Create achievement on platform","MODERATION ROLE CLAIM REQUIRED")]
   public async Task<ActionResult<ApiResponse<Achievement>>> CreateAchievement([FromForm] CreateAchievementRequest request)
   {
      var achievementCreationResult = await _achievementService.CreateAchievement(request);
      return Ok(ApiResponse<Achievement>.Ok("Successfully create achievement", achievementCreationResult));
   }
   
   [HttpPut]
   [SwaggerOperation("Update achievement","MODERATION ROLE CLAIM REQUIRED")]
   public async Task<ActionResult<ApiResponse<Achievement>>> UpdateAchievement([FromForm] UpdateAchievementRequest request)
   {
      var updateAchievementResult = await _achievementService.UpdateAchievement(request);
      return Ok(ApiResponse<Achievement>.Ok("Successfully updated achievement", updateAchievementResult));
   }

   [HttpDelete("{achievementId:guid}")]
   [SwaggerOperation("Delete achievement","MODERATION ROLE CLAIM REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteAchievement([FromRoute] Guid achievementId)
   {
      await _achievementService.DeleteAchievement(achievementId);
      return Ok(ApiResponse.Ok("Successfully deleted achievement"));
   }

   
   
}