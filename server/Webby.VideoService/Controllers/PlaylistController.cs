using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Helpers.Response;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Helpers.Jwt;
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
   public async Task<ActionResult<ApiResponse<List<PlaylistDto>>>> GetUserPlaylists([FromRoute] Guid userId)
   {
      var userPlaylists = await _playlistService.GetUserPlaylists(userId);
      return Ok(ApiResponse<List<PlaylistDto>>.Ok("Successfully retrieved user playlists",userPlaylists));
   }
   
   [HttpPost]
   [SwaggerOperation("Create user playlist","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PlaylistDto>>> CreatePlaylist([FromBody] CreatePlaylistRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      var playlistCreationResult = await _playlistService.CreatePlaylist(userId, request);
      return Ok(ApiResponse<PlaylistDto>.Ok("Successfully created playlist",playlistCreationResult));
   }

   [HttpDelete("{playlistId:guid}")]
   [SwaggerOperation("Delete user playlist", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteUserPlaylist([FromRoute] Guid playlistId)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      await _playlistService.DeletePlaylist(userId, playlistId);
      return Ok(ApiResponse.Ok("Successfully delete playlist"));
   }
   
   
}