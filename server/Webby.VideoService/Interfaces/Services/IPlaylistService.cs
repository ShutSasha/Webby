using Webby.VideoService.Dtos.Playlist;

namespace Webby.VideoService.Interfaces.Services;

public interface IPlaylistService
{
   Task<PlaylistDto> CreatePlaylist(Guid userId, CreatePlaylistRequest request);
   Task<List<PlaylistDto>> GetUserPlaylists(Guid userId);
   Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request);
   Task DeletePlaylist(Guid userId, Guid playlistId);
   Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId);
}