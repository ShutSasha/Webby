using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IPlaylistRepository : IRepository<Playlist>
{
   Task<List<Playlist>> GetUserPlaylists(Guid userId);
}