using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Models;

public class Notification
{
   public Guid NotificationId { get; set; }
   public Guid UserId { get; set; }
   public string Title { get; set; }
   public string Message { get; set; }
   public string TargetType { get; set; }
   public string TargetIdentifier { get; set; }
   public DateTime CreatedAt { get; set; }
   public NotificationStatus NotificationStatus { get; set; }
}