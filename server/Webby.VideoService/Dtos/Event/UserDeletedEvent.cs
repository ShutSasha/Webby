namespace Webby.VideoService.Dtos.Event;

public class UserDeletedEventDto
{
   public Guid UserId { get; set; }
   public DateTime? DeletedAt { get; set; }
}