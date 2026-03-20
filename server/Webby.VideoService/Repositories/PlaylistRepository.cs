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
   
   public async Task<Playlist?> FindByIdWithVideos(Guid playlistId)
   {
      return await _context.Playlists
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .FirstOrDefaultAsync(p => p.PlaylistId == playlistId);
   }

   public async Task<Playlist?> GetPlaylistDetails(Guid playlistId)
   {
      return await _context.Playlists
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .FirstOrDefaultAsync(p => p.PlaylistId == playlistId);

   }

   public async Task DeletePlaylistVideos(List<PlaylistVideo> videosToDelete)
   {
      _context.PlaylistVideos.RemoveRange(videosToDelete);
      await _context.SaveChangesAsync();
   }

   public async Task<(List<Playlist>, int)> GetPaginatedUserPlaylists(Guid userId, int page, int pageSize)
   {
      var query = _context.Playlists
         .Where(p => p.UserId == userId);

      var totalCount = await query.CountAsync();

      var playlists = await query
         .Skip((page - 1) * pageSize)
         .Take(pageSize)
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .ToListAsync();

      return (playlists, totalCount);
   }
}