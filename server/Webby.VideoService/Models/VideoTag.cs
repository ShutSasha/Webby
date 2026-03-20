namespace Webby.VideoService.Models;

public class VideoTag
{
   public Guid VideoId { get; set; }
   public Video Video { get; set; }
   public Guid TagId { get; set; }
   public Tag Tag { get; set; }
}