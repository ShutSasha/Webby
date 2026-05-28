using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Jwt;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Controllers;

[ApiController]
[Route("api/playlists")]
public class PlaylistController : ControllerBase
{
   private readonly IPlaylistService _playlistService;
   public PlaylistController(IPlaylistService playlistService)
   {
      _playlistService = playlistService;
   }
   
   [HttpGet("{userId:guid}")]
   [SwaggerOperation("Get user playlists")]
   public async Task<ActionResult<ApiResponse<PagedResponse<PlaylistPreviewDto>>>> GetUserPlaylists(
      [FromRoute] Guid userId,
      [FromQuery] GetUserPlaylistsRequest request,
      [FromQuery] string? videoId)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);

      var result = await _playlistService
         .GetUserPlaylists(
            requestUserId,
            videoId,
            userId,
            request);

      return Ok(ApiResponse<PagedResponse<PlaylistPreviewDto>>.Ok(
         "Successfully retrieved user playlists",
         result));
   }
   
   [HttpGet("{playlistId:guid}/videos/{videoId}/exists")]
   public async Task<ActionResult<ApiResponse<bool>>> CheckIfVideoExists([FromRoute] Guid playlistId, string videoId)
   {
      var checkVideoExistResult = await _playlistService.CheckIfVideoExistInPlaylist(playlistId, videoId);
      return Ok(ApiResponse<bool>.Ok("Successfully retrieve information",checkVideoExistResult));
   }
   
   [HttpGet("{playlistId:guid}/details")]
   [SwaggerOperation("Get playlist information")]
   public async Task<ActionResult<ApiResponse<GetPlaylistResponse>>> GetPlaylistInformation([FromRoute] Guid playlistId)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var playlistInformation = await _playlistService.GetPlaylistInformation(playlistId, requestUserId);
      return Ok(ApiResponse<GetPlaylistResponse>.Ok("Successfully retrieved playlist information", playlistInformation));
   }
   
   [HttpGet("search")]
   [SwaggerOperation("Search playlists route")]
   public async Task<ActionResult<ApiResponse<PagedResponse<SearchPlaylistDto>>>> SearchPlaylists(
      [FromQuery] SearchOptions searchOptions)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext, shouldThrowException: false);
      var searchPlaylistsResponse = await _playlistService.SearchPlaylists(requestUserId,searchOptions);
      return Ok(ApiResponse<PagedResponse<SearchPlaylistDto>>.Ok("Successfully retrieved public playlists",searchPlaylistsResponse));
   }
   
   [HttpPost]
   [SwaggerOperation("Create user playlist","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PlaylistDto>>> CreatePlaylist([FromBody] CreatePlaylistRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var playlistCreationResult = await _playlistService.CreatePlaylist(userId.Value, request);
      return Ok(ApiResponse<PlaylistDto>.Ok("Successfully created playlist",playlistCreationResult));
   }
   
   [HttpPost("videos")]
   [SwaggerOperation("Add or delete videos in playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> AddVideoToPlaylist([FromBody] AddVideoToPlaylistRequest request)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext)!;
      
      await _playlistService.AttachVideoToPlaylist(request.PlaylistId,request.VideoIds, requestUserId.Value);
      return Ok(ApiResponse.Ok("Successfully update playlist"));
   }

   [HttpPatch]
   [SwaggerOperation("Update playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PlaylistDto>>> UpdatePlaylist([FromBody] UpdatePlaylistRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var updatePlaylistResult = await _playlistService.UpdatePlaylist(userId.Value,request);
      return Ok(ApiResponse<PlaylistDto>.Ok("Successfully updated playlist", updatePlaylistResult));
   }
   
   [HttpDelete("{playlistId:guid}")]
   [SwaggerOperation("Delete user playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteUserPlaylist([FromRoute] Guid playlistId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      await _playlistService.DeletePlaylist(userId.Value, playlistId);
      return Ok(ApiResponse.Ok("Successfully delete playlist"));
   }
   
}