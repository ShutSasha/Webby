namespace Webby.NotificationService.Constants;

public static class WebSocketMethodNames
{
   public const string UpdateUnreadMessagesMethod = "UpdateUnreadNotificationsCount";
   public const string AuthErrorMethod = "AuthError";
   public const string ReceiveNotificationMethod = "ReceiveNotification";
}