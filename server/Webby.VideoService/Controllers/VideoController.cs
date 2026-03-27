using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Jwt;
using Webby.VideoService.Helpers.Response;
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

   [HttpGet("search")]
   [SwaggerOperation("Search video route")]
   public async Task<ActionResult<ApiResponse<PagedResponse<VideoDto>>>> SearchVideo([FromQuery] SearchOptions searchOptions)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var videos = await _videoService.SearchVideo(requestUserId,searchOptions);
      return Ok(ApiResponse<PagedResponse<VideoDto>>.Ok("Successfully retrieved video",videos));
   }
   
   [HttpGet("{playlistId:guid}/search")]
   [SwaggerOperation("Search video in playlist route")]
   public async Task<ActionResult<PagedResponse<VideoDto>>> SearchVideoInPlaylist(
      [FromRoute] Guid playlistId,
      [FromQuery] SearchOptions searchOptions)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var pagesResponse = await _videoService.SearchVideoInPlaylist(requestUserId,playlistId, searchOptions);
      return Ok(ApiResponse<PagedResponse<VideoDto>>.Ok("Successfully retrieved video from playlist",pagesResponse));
   }

   [HttpGet("users/{userId:guid}")]
   [SwaggerOperation("Get user videos")]
   public async Task<ActionResult<ApiResponse<PagedResponse<VideoDto>>>> GetUserVideos(
      [FromRoute] Guid userId,
      [FromQuery] int page = 1,
      [FromQuery] int pageSize = 10)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var result = await _videoService.GetUserVideos(userId, requestUserId, page, pageSize);
      return Ok(ApiResponse<PagedResponse<VideoDto>>.Ok("Successfully retrieved user videos", result));
   }
   
   //TODO: add recommendation videos route


   [HttpPost("upload")]
   [SwaggerOperation("Upload video file route", "AUTH REQUIRED")]
   [RequestSizeLimit(5L * 1024 * 1024 * 1024)]
   [RequestFormLimits(MultipartBodyLengthLimit = 5L * 1024 * 1024 * 1024)]
   public async Task<ActionResult<ApiResponse<UploadVideoResponse>>> UploadFileVideo([FromForm] UploadVideoRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var uploadResponse = await _videoService.UploadVideoFile(userId.Value, request);
      return Ok(ApiResponse<UploadVideoResponse>.Ok("Successfully prepared for uploading video file",uploadResponse));
   }
   
   [HttpPost]
   [SwaggerOperation("Add information to video", "AUTH REQUIRED")]
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

   [HttpDelete("{videoId:guid}/cancel")]
   [SwaggerOperation("Cancel uploading video route")]
   public async Task<ActionResult<ApiResponse>> CancelVideo([FromRoute] Guid videoId)
   {
      await _videoService.CancelVideoUploading(videoId);
      return Ok(ApiResponse.Ok("Successfully cancel video upload"));
   }
   
}