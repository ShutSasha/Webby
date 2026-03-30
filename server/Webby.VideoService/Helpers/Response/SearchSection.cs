namespace Webby.VideoService.Helpers.Response;

public class SearchSection<T> where T : class
{
   public List<T> Items { get; set; } = [];
   public int Page { get; set; } = 1;
   public string? NextPageToken { get; set; }
   public int TotalCount { get; set; } = 0; 
   public int RemainingCount => Math.Max(0, TotalCount - Items.Count);
}