using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Services;

public interface IPlaylistService
{
   Task<Playlist> GetPlaylistById(Guid playlistId);
   Task<PlaylistDto> CreatePlaylist(Guid userId, CreatePlaylistRequest request);
   Task<PagedResponse<PlaylistPreviewDto>> GetUserPlaylists(Guid? requestUserId, string? videoId, Guid userId, GetUserPlaylistsRequest request);
   Task<PlaylistDto> UpdatePlaylist(Guid requestUserId, UpdatePlaylistRequest request);
   Task DeletePlaylist(Guid userId, Guid playlistId);
   Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId, Guid? requestedUserId);
   Task AttachVideoToPlaylist(Guid playlistId, List<string> videoIds, Guid requestUserId);
   Task AttachVideoToPlaylists(List<Guid> playlistIds, List<string> videoIds, Guid requestUserId);
   Task AttachStreamToPlaylist(Guid playlistId, List<string> streamIds, Guid requestUserId);
   Task AttachStreamToPlaylists(List<Guid> playlistIds, List<string> videoIds, Guid requestUserId);
   Task<PagedResponse<SearchPlaylistDto>> SearchPlaylists(Guid? requestUserId, SearchOptions searchOptions);
   Task<bool> CheckIfVideoExistInPlaylist(Guid playlistId, string videoId);

}