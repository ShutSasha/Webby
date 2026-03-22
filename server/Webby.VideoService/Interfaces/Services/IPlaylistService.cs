using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IPlaylistService
{
   Task<PlaylistDto> CreatePlaylist(Guid userId, CreatePlaylistRequest request);
   Task<PagedResponse<PlaylistDto>> GetUserPlaylists(Guid? requestUserId, Guid userId, int page, int pageSize);
   Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request);
   Task DeletePlaylist(Guid userId, Guid playlistId);
   Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId);
   Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<Guid> videoIds);
}