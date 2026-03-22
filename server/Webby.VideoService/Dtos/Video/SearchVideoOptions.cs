namespace Webby.VideoService.Dtos.Video;

public class SearchVideoOptions
{
   public string? SearchText { get; set; }
   public int Page { get; set; } = 1;
   public int PageSize { get; set; } = 10;
}