using Webby.VideoService.Dtos.User;

namespace Webby.VideoService.Dtos.Video;

public class VideoDto
{
   public Guid VideoId { get; set; }
   public string Name { get; set; }
   public int Views { get; set; }
   public DateTime CreatedAt { get; set; }
   public UserVideoDto User { get; set; }
}