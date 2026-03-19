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


   [HttpGet("{videoId:guid}")]
   [SwaggerOperation("Get video information")]
   public async Task<ActionResult<ApiResponse<GetVideoInformationResponse>>> GetVideoInformation(
      [FromRoute] Guid videoId)
   {
      var requestedUserId = JwtHelper.ExtractUserId(HttpContext,false);
      var getVideoInformationResult = await _videoService.GetVideoInformation(videoId, requestedUserId);

      return Ok(ApiResponse<GetVideoInformationResponse>.Ok("Successfully retrieved video information",
         getVideoInformationResult));
   }

   [HttpGet("users/{userId:guid}")]
   [SwaggerOperation("Get user videos")]
   public async Task<ActionResult<ApiResponse<List<VideoDto>>>> GetUserVideos([FromRoute] Guid userId)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var userVideos = await _videoService.GetUserVideos(userId,requestUserId);
      return Ok(ApiResponse<List<VideoDto>>.Ok("Successfully retrieved user videos", userVideos));
   }
   
   
   [HttpPost]
   [SwaggerOperation("Create video route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> CreateVideo([FromForm] CreateVideoRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      await _videoService.CreateVideo(userId.Value, request);
      return Ok(ApiResponse.Ok("Successfully create video"));
   }

   [HttpPatch]
   [SwaggerOperation("Update video information route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> UpdateVideo([FromForm] UpdateVideoRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      await _videoService.UpdateVideoInformation(userId.Value,request);
      return Ok(ApiResponse.Ok("Successfully update video"));
   }

   [HttpDelete("{videoId:guid}")]
   [SwaggerOperation("Delete video route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteVideo([FromRoute] Guid videoId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      await _videoService.DeleteVideo(userId.Value, videoId);
      return Ok(ApiResponse.Ok("Successfully delete video"));
   }
   
}