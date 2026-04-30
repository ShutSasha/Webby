using Webby.VideoService.Dtos.User;

namespace Webby.VideoService.Dtos.Video;

public class PreviewVideoDto
{
   public string VideoId { get; init; }
   public string Name { get; init; }
   public string PreviewUrl { get; init; }
   public int Views { get; init; }
   public bool IsPrivate { get; init; }
   public DateTime CreatedAt { get; init; }
   public long Duration { get; init; }
   public UserVideoDto User { get; set; }
}