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
   
   public async Task<List<Playlist>> GetPlaylistsDetails(List<Guid> playlistIds)
   {
      return await _context.Playlists
         .Where(p => playlistIds.Contains(p.PlaylistId))
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .ToListAsync();
   }
   
   public async Task DeletePlaylistVideos(List<PlaylistVideo> videosToDelete)
   {
      _context.PlaylistVideos.RemoveRange(videosToDelete);
      await _context.SaveChangesAsync();
   }
   
   public async Task<(List<Playlist>, int)> GetPaginatedUserPlaylists(bool shouldShowPrivate, Guid userId, int page, int pageSize)
   {
      var query = _context.Playlists
         .Where(p => p.UserId == userId && (shouldShowPrivate || !p.IsPrivate));
   
      var totalCount = await query.CountAsync();
   
      var playlists = await query
         .Skip((page - 1) * pageSize)
         .Take(pageSize)
         .Include(p => p.PlaylistVideos)
         .ThenInclude(pv => pv.Video)
         .ToListAsync();
   
      return (playlists, totalCount);
   }
   
   public async Task<(List<Playlist> Items, int Total)> SearchPlaylistsAsync(
      string? searchText,
      int skip,
      int take)
   {
      var baseQuery = _context.Playlists
         .Where(p => !p.IsPrivate);
   
      if (string.IsNullOrWhiteSpace(searchText))
      {
         var query = baseQuery.OrderByDescending(p => p.CreatedAt);
   
         var count = await query.CountAsync();
         var data = await query.Skip(skip).Take(take).ToListAsync();
   
         return (data, count);
      }
   
      var search = searchText.Trim();
      var likePattern = $"%{search}%";
   
      if (search.Length < 3)
      {
         var query = baseQuery
            .Where(p => p.Name.Contains(search))
            .OrderByDescending(p => p.CreatedAt);
   
         var count = await query.CountAsync();
         var data = await query.Skip(skip).Take(take).ToListAsync();
   
         return (data, count);
      }
   
      var filter = @"
        FROM ""Playlists""
        WHERE 
            (""Name"" <% {0} OR ""Name"" ILIKE {1})
            AND ""IsPrivate"" = FALSE
    ";
   
      var total = await _context.Playlists
         .FromSqlRaw($"SELECT * {filter}", search, likePattern)
         .CountAsync();
   
      var items = await _context.Playlists
         .FromSqlRaw($@"
            SELECT *
            {filter}
            ORDER BY 
                (CASE WHEN ""Name"" ILIKE {{1}} THEN 1 ELSE 0 END) DESC,
                word_similarity({{0}}, ""Name"") DESC
            LIMIT {{2}}
            OFFSET {{3}}
        ", search, likePattern, take, skip)
         .ToListAsync();
   
      return (items, total);
   }
   
   public async Task<bool> CheckIsVideoAdded(string videoId, Guid playlistId)
   {
      if (Guid.TryParse(videoId, out var localGuid))
      {
         return await _context.PlaylistVideos
            .AnyAsync(pv => pv.PlaylistId == playlistId &&
                            (pv.VideoId == localGuid || pv.ExternalVideoId == videoId));
      }
   
      return await _context.PlaylistVideos
         .AnyAsync(pv => pv.PlaylistId == playlistId && pv.ExternalVideoId == videoId);
   }
   
   public async Task<HashSet<Guid>> GetPlaylistIdsContainingVideo(
      string videoId,
      List<Guid> playlistIds)
   {
      if (string.IsNullOrEmpty(videoId) || playlistIds == null || !playlistIds.Any())
         return new HashSet<Guid>();
   
      List<Guid> ids;
      
      if (Guid.TryParse(videoId, out var localGuid))
      {
         ids = await _context.PlaylistVideos
            .Where(pv => playlistIds.Contains(pv.PlaylistId) && 
                         (pv.VideoId == localGuid || pv.ExternalVideoId == videoId))
            .Select(pv => pv.PlaylistId)
            .ToListAsync();
      }
      else
      {
         ids = await _context.PlaylistVideos
            .Where(pv => playlistIds.Contains(pv.PlaylistId) && pv.ExternalVideoId == videoId)
            .Select(pv => pv.PlaylistId)
            .ToListAsync();
      }
   
      return ids.ToHashSet();
   }
}