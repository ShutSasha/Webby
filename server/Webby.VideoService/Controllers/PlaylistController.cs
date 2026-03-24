using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.VideoService.Dtos.Playlist;
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
   public async Task<ActionResult<ApiResponse<PagedResponse<PlaylistDto>>>> GetUserPlaylists(
      [FromRoute] Guid userId,
      [FromQuery] int page = 1,
      [FromQuery] int pageSize = 10)
   {
      var result = await _playlistService.GetUserPlaylists(userId, page, pageSize);

      return Ok(ApiResponse<PagedResponse<PlaylistDto>>.Ok(
         "Successfully retrieved user playlists",
         result));
   }

   [HttpGet("{playlistId:guid}/details")]
   [SwaggerOperation("Get playlist information")]
   public async Task<ActionResult<ApiResponse<GetPlaylistResponse>>> GetPlaylistInformation([FromRoute] Guid playlistId)
   {
      var playlistInformation = await _playlistService.GetPlaylistInformation(playlistId);
      return Ok(ApiResponse<GetPlaylistResponse>.Ok("Successfully retrieved playlist information", playlistInformation));
   }
   
   [HttpPost]
   [SwaggerOperation("Create user playlist","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PlaylistDto>>> CreatePlaylist([FromBody] CreatePlaylistRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var playlistCreationResult = await _playlistService.CreatePlaylist(userId.Value, request);
      return Ok(ApiResponse<PlaylistDto>.Ok("Successfully created playlist",playlistCreationResult));
   }
   
   //TODO: Measure response time in stress testing
   [HttpPost("videos")]
   [SwaggerOperation("Add or delete videos in playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> AddVideoToPlaylist([FromBody] AddVideoToPlaylistRequest request)
   {
      await _playlistService.AttachVideoToPlaylist(request.PlaylistId,request.VideoIds);
      return Ok(ApiResponse.Ok("Successfully added videos to playlist"));
   }

   [HttpPatch]
   [SwaggerOperation("Update playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PlaylistDto>>> UpdatePlaylist([FromBody] UpdatePlaylistRequest request)
   {
      var updatePlaylistResult = await _playlistService.UpdatePlaylist(request);
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