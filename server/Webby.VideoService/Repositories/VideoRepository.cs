using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;

namespace Webby.VideoService.Repositories;

public class VideoRepository : GenericRepository<Video>,IVideoRepository
{
   public VideoRepository(AppDbContext context) : base(context)
   {
   }
}