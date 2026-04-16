using Webby.VideoService.Dtos.User;

namespace Webby.VideoService.Dtos.Stream;

public class StreamDto
{
   public string StreamId { get; set; }
   public string Name { get; set; }
   public string Source { get; set; }
   public string PreviewUrl { get; set; }
   public DateTime StartedAt { get; set; }
   public int Viewers { get; set; }
   public string StreamUrl { get; set; }
   public UserVideoDto? User { get; set; }
}