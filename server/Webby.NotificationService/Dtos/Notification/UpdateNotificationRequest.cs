namespace Webby.NotificationService.Dtos.Notification;

public class UpdateNotificationRequest
{
   public Guid NotificationId { get; set; }
   public string Title { get; set; }
   public string Message { get; set; }
   public string TargetType { get; set; }
   public string TargetIdentifier { get; set; }
}