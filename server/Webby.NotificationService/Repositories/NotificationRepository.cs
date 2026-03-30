using Microsoft.EntityFrameworkCore;
using Webby.NotificationService.Data;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Models;
using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Repositories;

public class NotificationRepository : GenericRepository<Notification>, INotificationRepository
{
   public NotificationRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<int> CountUnreadMessages(Guid userId)
   {
      return await _context.Notifications
         .Where(n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Unread)
         .CountAsync();
   }

   public async Task ChangeReadStatus(List<Guid> ids)
   {
      if (ids ==  null || !ids.Any())
      {
         return;
      }

      await _context.Notifications
         .Where(n => ids.Contains(n.NotificationId) && n.NotificationStatus == NotificationStatus.Unread)
         .ExecuteUpdateAsync(s => s.SetProperty(n => n.NotificationStatus, NotificationStatus.Read));
   }
}