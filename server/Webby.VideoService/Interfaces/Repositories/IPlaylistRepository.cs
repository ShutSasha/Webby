using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IPlaylistRepository : IRepository<Playlist>
{
   Task<List<Playlist>> GetUserPlaylists(Guid userId);
   Task AddPlaylistVideos(List<PlaylistVideo> playlistVideos);
   Task<Playlist?> FindByIdWithVideos(Guid playlistId);
   Task<Playlist?> GetPlaylistDetails(Guid playlistId);
   Task DeletePlaylistVideos(List<PlaylistVideo> videosToDelete);
   Task<(List<Playlist>, int)> GetPaginatedUserPlaylists(bool shouldShowPrivate, Guid userId, int page, int pageSize);
}