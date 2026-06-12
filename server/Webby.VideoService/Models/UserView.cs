namespace Webby.VideoService.Models;

public class UserView
{
   public Guid VideoId { get; set; }
   public Guid UserId { get; set; }
   public DateTime WatchedAt { get; set; } = DateTime.UtcNow;
}