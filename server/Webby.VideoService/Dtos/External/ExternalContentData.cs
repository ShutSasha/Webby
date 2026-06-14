using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.Video;

namespace Webby.VideoService.Dtos.External;

public class ExternalContentData
{
   public Dictionary<string, VideoDto> YoutubeVideos { get; set; } = new();
   public Dictionary<string, StreamDto> TwitchStreams { get; set; } = new();
   public int UnavailableCount { get; set; }
}