namespace Webby.NotificationService.Dtos.Notification;

public class GetNotificationsCountResponse
{
   public int CountOfUnreadMessages { get; set; }
   public int CountOfReadMessages { get; set; }
}