using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.Stream.Enums;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IStreamService
{
   Task<PagedResponse<StreamDto>> SearchStream(SearchStreamOptions searchOptions);
   Task<StreamDto> GetStreamById(string streamerId);
}