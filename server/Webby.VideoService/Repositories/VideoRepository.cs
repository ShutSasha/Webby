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
}