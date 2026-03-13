using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Dtos;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Helpers.Jwt;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;

namespace Webby.UserService.Controllers;

[ApiController]
[Route("api/users")]
public class UserController : ControllerBase
{
   private readonly IUserService _userService;
   
   public UserController(IUserService userService)
   {
      _userService = userService;
   }

   [HttpGet("{userId:guid}")]
   [SwaggerOperation("Get user information by id")]
   public async Task<ActionResult<ApiResponse<UserDto>>> GetUserInformation([FromRoute] Guid userId)
   {
      var userExtractionResult = await _userService.GetUserInformation(userId);
      return Ok(ApiResponse.Ok("Successfully extract user", userExtractionResult));
   }
   
   [HttpGet("{userId:guid}/follows")]
   [SwaggerOperation("Get user follows")]
   public async Task<ActionResult<ApiResponse<List<UserFollowersDto>>>> GetUserFollows([FromRoute] Guid userId)
   {
      var userFollows = await _userService.GetUserFollows(userId);
      return Ok(ApiResponse<List<UserFollowersDto>>.Ok("Successfully retrieved user follows", userFollows));
   }
   
   [HttpGet("{userId:guid}/followers")]
   [SwaggerOperation("Get user followers")]
   public async Task<ActionResult<ApiResponse<List<UserFollowersDto>>>> GetUserFollowers([FromRoute] Guid userId)
   {
      var userFollowers = await _userService.GetUserFollowers(userId);
      return Ok(ApiResponse<List<UserFollowersDto>>.Ok("Successfully retrieved user followers", userFollowers));
   }
   
   [HttpPost("{followId:guid}/follow")]
   [SwaggerOperation("Follow or unfollow user","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> Follow(Guid followId)
   {
      var followerId = JwtHelper.ExtractUserId(HttpContext);

      var resultMessage = await _userService.ProcessFollow(new UserFollowRequest
      {
         UserId = followId,
         FollowerId = followerId
      });

      return Ok(ApiResponse.Ok(resultMessage));
   }
   
   [HttpPost("achievements/unlock")]
   [SwaggerOperation("Test unlock achievement to user")]
   public async Task<ActionResult<ApiResponse>> UnlockUserAchievement([FromBody] UnlockUserAchievementRequest request)
   {
      await _userService.UnlockAchievement(request.UserId, request.AchievementId);
      return Ok(ApiResponse.Ok("Successfully unlock user achievement"));
   }

   [HttpPost("achievements/{achievementId:guid}")]
   [SwaggerOperation("Pin user achievement","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> PinUserAchievement([FromRoute] Guid achievementId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      await _userService.PinUserAchievement(userId, achievementId);
      return Ok(ApiResponse.Ok("Successfully pinned user achievement"));

   }
   
   [HttpPatch]
   [SwaggerOperation("Update user text information","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<UserProfileResponse>>> UpdateUserInformation([FromBody] UpdateUserRequest request)
   {
      var userUpdateResult = await _userService.UpdateUserInformation(request);
      return Ok(ApiResponse<UserProfileResponse>.Ok("Successfully update user", userUpdateResult));
   }

   [HttpPatch("update-user-avatar/{userId:guid}")]
   [SwaggerOperation("Update user icon","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<UserDto>>> UpdateUserAvatar(IFormFile file, [FromRoute] Guid userId)
   {
      if (file.Length == 0)
      {
         throw new ApiException("Update user avatar error", 400, "file is empty");
      }
      
      await using var fileStream = file.OpenReadStream();
      var fileName = file.FileName;
      var contentType = file.ContentType;
      var updateUserIconResult = await _userService.EditUserIcon(userId, fileName, fileStream, contentType);

      return Ok(ApiResponse<UserDto>.Ok("Successfully update user icon", updateUserIconResult));
   }
   
   
   [HttpDelete("achievements/{achievementId:guid}")]
   [SwaggerOperation("Unpin user achievement","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> UnpinUserAchievement([FromRoute] Guid achievementId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      await _userService.UnpinUserAchievement(userId, achievementId);
      return Ok(ApiResponse.Ok("Successfully unpinned user achievement"));
   }
   
}