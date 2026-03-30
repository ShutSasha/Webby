using Webby.NotificationService.Dtos.Notification;
using Webby.NotificationService.Helpers.Exception;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Interfaces.Services;
using Webby.NotificationService.Models;
using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Services;

public class NotificationService : INotificationService
{
   private readonly INotificationRepository _notificationRepository;

   public NotificationService(INotificationRepository notificationRepository)
   {
      _notificationRepository = notificationRepository;
   }
   
   public async Task<Notification> CreateNotification(CreateNotificationRequest request)
   {
      var notification = new Notification()
      {
         NotificationId = Guid.NewGuid(),
         UserId = request.UserId,
         Title = request.Title,
         Message = request.Message,
         CreatedAt = DateTime.UtcNow,
         NotificationStatus = NotificationStatus.Unread,
      };
      await _notificationRepository.Add(notification);
      return notification;
   }

   public async Task<Notification> UpdateNotification(UpdateNotificationRequest request)
   {
      var notification = await _notificationRepository.FindById(request.NotificationId)
                         ?? throw new ApiException("Update notification error", 404, "Notification wasn't found");

      notification.Message = request.Message;
      notification.Title = request.Title;

      await _notificationRepository.Update(notification);

      return notification;
   }

   public async Task<List<Notification>> GetUnreadNotifications(Guid userId)
   {
      return (await _notificationRepository
         .GetByPredicate(n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Unread)).ToList();
   }

   public async Task<List<Notification>> GetReadNotifications(Guid userId)
   {
      return (await _notificationRepository
         .GetByPredicate(n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Read)).ToList();
   }

   public async Task<int> GetNotificationsCount(Guid userId)
      => await _notificationRepository.CountUnreadMessages(userId);

   public async Task DeleteNotification(Guid requestUserId, Guid notificationId)
   {
      var notification = await _notificationRepository.FindById(notificationId)
                         ?? throw new ApiException("Delete notification error", 404, "Notification wasn't found");
      if (notification.UserId != requestUserId)
      {
         throw new ApiException("Delete notification error", 403, "You can't delete this notification");
      }
      await _notificationRepository.DeleteAsync(notificationId);
   }

   public async Task ChangeReadStatus(List<Guid> notificationIds)
      => await _notificationRepository.ChangeReadStatus(notificationIds);
}
   
