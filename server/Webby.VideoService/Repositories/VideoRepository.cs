using Microsoft.EntityFrameworkCore;
using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;

namespace Webby.VideoService.Repositories;

public class VideoRepository : GenericRepository<Video>,IVideoRepository
{
   public VideoRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<Video> GetVideoInformationById(Guid videoId)
   {
      return await _context.Videos
         .Include(v => v.VideoTags)
         .FirstAsync(v => v.VideoId == videoId);
   }
   public async Task<(List<Video>, int)> GetPaginatedUserVideos(Guid userId, bool isOwner, int page, int pageSize)
   {
      var query = _context.Videos
         .Where(v => v.UserId == userId && (isOwner || !v.IsPrivate));

      var totalCount = await query.CountAsync();

      var videos = await query
         .OrderByDescending(v => v.CreatedAt)
         .Skip((page - 1) * pageSize)
         .Take(pageSize)
         .ToListAsync();

      return (videos, totalCount);
   }
   
   public async Task<(List<Video> Items, int Total)> SearchAsync(
      string? searchText,
      int skip,
      int take)
   {
      if (string.IsNullOrWhiteSpace(searchText))
      {
         var query = _context.Videos.OrderByDescending(v => v.CreatedAt);
         var count = await query.CountAsync();
         var data = await query.Skip(skip).Take(take).ToListAsync();
         return (data, count);
      }

      var search = searchText.Trim();
      var likePattern = $"%{search}%";

      if (search.Length < 3)
      {
         var query = _context.Videos
            .Where(v => v.Name.Contains(search))
            .OrderByDescending(v => v.CreatedAt);
         var count = await query.CountAsync();
         var data = await query.Skip(skip).Take(take).ToListAsync();
         return (data, count);
      }
      
      var total = await _context.Videos
         .FromSqlInterpolated($@"
        SELECT * FROM ""Videos"" 
        WHERE ""Name"" <% {search} OR ""Name"" ILIKE {likePattern}
    ")
         .CountAsync();

      var items = await _context.Videos
         .FromSqlInterpolated($@"
         SELECT *
         FROM ""Videos""
         WHERE ""Name"" <% {search} OR ""Name"" ILIKE {likePattern}
         ORDER BY 
            (CASE WHEN ""Name"" ILIKE {likePattern} THEN 1 ELSE 0 END) DESC,
            word_similarity({search}, ""Name"") DESC
         LIMIT {take}
         OFFSET {skip}
      ")
         .ToListAsync();

      return (items, total);
   }
}