namespace Webby.VideoService.Helpers.Response;

public class PagedResponse<T> where T : class
{
   public List<T> Items { get; set; } = [];
   public int Page { get; set; }
   public int PageSize { get; set; }
   public int TotalCount { get; set; }
}