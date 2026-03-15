using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services;

public class VideoService : IVideoService
{
   private readonly IVideoRepository _videoRepository;
   public VideoService(IVideoRepository videoRepository)
   {
      _videoRepository = videoRepository;
   }
   
   
   
}