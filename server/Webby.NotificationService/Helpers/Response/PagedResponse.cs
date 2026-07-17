namespace Webby.NotificationService.Helpers.Response;

public class PagedResponse<T> where T : class
{
   public List<T> Items { get; set; } = [];
   public int Page { get; set; } = 1;
   public int PageSize { get; set; } = 10;
   public int TotalCount { get; set; } = 0;
}