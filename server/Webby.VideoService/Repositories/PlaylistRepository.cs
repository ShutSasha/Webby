using Microsoft.EntityFrameworkCore;
using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;

namespace Webby.VideoService.Repositories;

public class PlaylistRepository : GenericRepository<Playlist>, IPlaylistRepository
{
   public PlaylistRepository(AppDbContext context) : base(context)
   {
   }
   
   public async Task<List<Playlist>> GetUserPlaylists(Guid userId)
   {
      return await _context.Playlists
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .Where(p => p.UserId == userId)
         .ToListAsync();
   }

   public async Task AddPlaylistVideos(List<PlaylistVideo> playlistVideos)
   {
      await _context.PlaylistVideos.AddRangeAsync(playlistVideos);
      await _context.SaveChangesAsync();
   }
}