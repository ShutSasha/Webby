using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Dtos.Notification;

public class CreateNotificationRequest
{
   public Guid UserId { get; set; }
   public string Title { get; set; }
   public string Message { get; set; }
   public NotificationTargetType TargetType { get; set; }
   public string  TargetIdentifier { get; set; }
}