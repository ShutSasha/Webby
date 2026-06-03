using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IPlaylistRepository : IRepository<Playlist>
{
   Task<List<Playlist>> GetUserPlaylists(Guid userId);
   Task AddPlaylistVideos(List<PlaylistVideo> playlistVideos);
   Task<Playlist?> FindByIdWithVideos(Guid playlistId);
   Task<Playlist?> GetPlaylistDetails(Guid playlistId);
   Task<List<Playlist>> GetPlaylistsDetails(List<Guid> playlistIds);
   Task DeletePlaylistVideos(List<PlaylistVideo> videosToDelete);
   Task<(List<Playlist>, int)> GetPaginatedUserPlaylists(bool shouldShowPrivate, Guid userId, int page, int pageSize);
   Task<(List<Playlist> Items, int Total)> SearchPlaylistsAsync(string? searchText, int skip, int take);
   Task<bool> CheckIsVideoAdded(string videoId, Guid playlistId);
   Task<HashSet<Guid>> GetPlaylistIdsContainingVideo(string videoId, List<Guid> playlistIds);
   Task DeleteUserPlaylists(Guid userId);
}