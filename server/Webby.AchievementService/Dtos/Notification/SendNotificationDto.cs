namespace Webby.AchievementService.Dtos.Notification;

public class SendNotificationDto
{
   public string UserId { get; set; }
   public string Title { get; set; }
   public string Message { get; set; }
   public string TargetIdentifier { get; set; }
}