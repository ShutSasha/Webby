using Webby.VideoService.Dtos.Stream;

namespace Webby.VideoService.Interfaces.Services;

public interface ITwitchSearchService : IExternalVideoSearchService<StreamDto>
{
   Task<List<StreamDto>> GetList(List<string> sourceIds);
}