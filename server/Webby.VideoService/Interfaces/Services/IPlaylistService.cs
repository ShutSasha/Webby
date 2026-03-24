using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IPlaylistService
{
   Task<PlaylistDto> CreatePlaylist(Guid userId, CreatePlaylistRequest request);
   Task<PagedResponse<PlaylistPreviewDto>> GetUserPlaylists(Guid? requestUserId, Guid? videoId, Guid userId, SearchOptions searchOptions);
   Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request);
   Task DeletePlaylist(Guid userId, Guid playlistId);
   Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId, Guid? requestedUserId);
   Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<Guid> videoIds);
   Task<PagedResponse<PlaylistDto>> SearchPlaylists(Guid? requestUserId, SearchOptions searchOptions);
}