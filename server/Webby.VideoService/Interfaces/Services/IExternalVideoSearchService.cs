using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IExternalVideoSearchService
{
   Task<PagedResponse<VideoDto>> SearchAsync(SearchVideoOptions options);
   Task<VideoDto> FindById(string videoId);
}

