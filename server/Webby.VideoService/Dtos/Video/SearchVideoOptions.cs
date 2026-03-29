using Webby.VideoService.Dtos.Playlist;

namespace Webby.VideoService.Dtos.Video;

public class SearchVideoOptions : SearchOptions
{
   public SearchPlatforms SearchPlatform { get; set; }
   public string? NextPageToken { get; set; }
}