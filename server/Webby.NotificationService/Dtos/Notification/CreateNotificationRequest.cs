namespace Webby.NotificationService.Dtos.Notification;

public class CreateNotificationRequest
{
   public Guid UserId { get; set; }
   public string Title { get; set; }
   public string Message { get; set; }
}