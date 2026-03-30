using Webby.NotificationService.Dtos.Notification;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Models;

namespace Webby.NotificationService.Interfaces.Services;

public interface INotificationService
{
   Task<Notification> CreateNotification(CreateNotificationRequest request);
   Task<Notification> UpdateNotification(UpdateNotificationRequest request);
   Task<List<Notification>> GetUnreadNotifications(Guid userId);
   Task<List<Notification>> GetReadNotifications(Guid userId);
   Task<int> GetNotificationsCount(Guid userId);
   Task DeleteNotification(Guid requestUserId, Guid notificationId);
   Task ChangeReadStatus(List<Guid> notificationIds);
}