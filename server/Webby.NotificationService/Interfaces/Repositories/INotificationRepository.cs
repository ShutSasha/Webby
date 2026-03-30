using Webby.NotificationService.Models;

namespace Webby.NotificationService.Interfaces.Repositories;

public interface INotificationRepository : IRepository<Notification>
{
   Task<int> CountUnreadMessages(Guid userId);
   Task ChangeReadStatus(List<Guid> ids);
}