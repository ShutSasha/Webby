using Webby.NotificationService.Models;
using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Interfaces.Repositories;

public interface INotificationRepository : IRepository<Notification>
{
   Task<int> CountNotifications (Guid userId,NotificationStatus status);
   Task ChangeReadStatus(List<Guid> ids);
}