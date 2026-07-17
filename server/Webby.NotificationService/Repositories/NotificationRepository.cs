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

   public async Task<int> CountNotifications(Guid userId,NotificationStatus status)
   {
      return await _context.Notifications
         .Where(n => n.UserId == userId && n.NotificationStatus == status)
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

   public async Task<bool> IsRecentDuplicateAsync(Guid userId, NotificationTargetType targetType, string targetIdentifier, TimeSpan timeWindow)
   {
      var thresholdTime = DateTime.UtcNow.Subtract(timeWindow);

      return await _context.Notifications.AnyAsync(n => 
         n.UserId == userId && 
         n.TargetType == targetType && 
         n.TargetIdentifier == targetIdentifier && 
         n.CreatedAt >= thresholdTime);
   }
}