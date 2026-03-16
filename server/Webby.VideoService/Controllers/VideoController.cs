using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Helpers.Response;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Jwt;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Controllers;

[ApiController]
[Route("api/videos")]
public class VideoController: ControllerBase
{
   private readonly IVideoService _videoService;
   public VideoController(IVideoService videoService)
   {
      _videoService = videoService;
   }

   [HttpPost]
   [SwaggerOperation("Create video route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> CreateVideo([FromForm] CreateVideoRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      await _videoService.CreateVideo(userId, request);
      return Ok(ApiResponse.Ok("Successfully create video"));
   }

   [HttpDelete("{videoId:guid}")]
   [SwaggerOperation("Delete video route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteVideo([FromRoute] Guid videoId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      await _videoService.DeleteVideo(userId, videoId);
      return Ok(ApiResponse.Ok("Successfully delete video"));
   }
   
}