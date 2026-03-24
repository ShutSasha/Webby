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
   
}