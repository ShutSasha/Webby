namespace Webby.VideoService.Dtos.External;

public class YouTubeChannelResponse
{
   public List<ChannelItem> Items { get; set; }
}

public class ChannelItem
{
   public string Id { get; set; }
   public ChannelSnippet Snippet { get; set; }
}

public class ChannelSnippet
{
   public YouTubeVideoResponse.Thumbnails Thumbnails { get; set; }
}