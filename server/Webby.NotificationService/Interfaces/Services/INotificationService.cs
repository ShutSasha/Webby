using Webby.NotificationService.Dtos.Notification;
using Webby.NotificationService.Dtos.Pagination;
using Webby.NotificationService.Helpers.Response;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Models;

namespace Webby.NotificationService.Interfaces.Services;

public interface INotificationService
{
   Task<Notification> CreateNotification(CreateNotificationRequest request);
   Task<Notification> UpdateNotification(UpdateNotificationRequest request);
   Task<PagedResponse<Notification>> GetUnreadNotifications(Guid userId, PaginationRequest request);
   Task<PagedResponse<Notification>> GetReadNotifications(Guid userId, PaginationRequest request);
   Task<int> GetNotificationsCount(Guid? userId);
   Task DeleteNotification(Guid requestUserId, Guid notificationId);
   Task ChangeReadStatus(List<Guid> notificationIds);
   Task<GetNotificationsCountResponse> GetUsersNotificationsCount(Guid userId);
   Task<PagedResponse<Notification>> GetUserNotifications(Guid userId, PaginationRequest request);
}