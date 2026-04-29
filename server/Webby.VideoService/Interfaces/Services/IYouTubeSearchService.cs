using Webby.VideoService.Dtos.Video;

namespace Webby.VideoService.Interfaces.Services;

public interface IYouTubeSearchService : IExternalVideoSearchService<VideoDto>
{
   Task<List<VideoDto>> GetList(List<string> sourceIds);
}