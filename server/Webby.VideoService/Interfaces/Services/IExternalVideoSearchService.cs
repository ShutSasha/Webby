using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IExternalVideoSearchService<TSource> where TSource : class
{
   Task<PagedResponse<TSource>> SearchAsync(string? searchText, int pageSize, int page, string? nextPageToken);
   Task<TSource> FindById(string sourceId);
   Task<List<TSource>> GetList(List<string> sourceIds);
}

