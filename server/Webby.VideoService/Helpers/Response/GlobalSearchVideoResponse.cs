using Webby.VideoService.Dtos.Video;

namespace Webby.VideoService.Helpers.Response;

public class GlobalSearchVideoResponse
{
   public SearchSection<VideoDto> WebbyVideos { get; set; } = new();
   public SearchSection<VideoDto> YouTubeVideos { get; set; } = new();
}

