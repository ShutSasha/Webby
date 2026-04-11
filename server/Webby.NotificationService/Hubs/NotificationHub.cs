using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.SignalR;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Hubs;

public class NotificationHub : Hub
{
   private readonly INotificationRepository _notificationRepository;

   public NotificationHub(INotificationRepository notificationRepository)
   {
      _notificationRepository = notificationRepository;
   }
   
   public override async Task OnConnectedAsync()
   {
      if (Context.User?.Identity?.IsAuthenticated != true)
      {
         await Clients.Caller.SendAsync("AuthError", "Unauthorized: Token is missing or expired");
         Context.Abort();
         return;
      }

      var userIdString = Context.UserIdentifier;

      if (Guid.TryParse(userIdString, out Guid userId))
      {
         var unreadNotificationCount = await _notificationRepository.CountNotifications(userId, NotificationStatus.Unread);
         await Clients.Caller.SendAsync("UpdateUnreadCount", unreadNotificationCount);
      }
      
      await base.OnConnectedAsync();
   }
   
   public override async Task OnDisconnectedAsync(Exception? exception)
   {
      if (exception != null)
      {
         Console.WriteLine($"DISCONNECTED ERROR ({Context.UserIdentifier}): {exception.Message}");
         Console.WriteLine($"StackTrace: {exception.StackTrace}");
      }

      await base.OnDisconnectedAsync(exception);
   }
}