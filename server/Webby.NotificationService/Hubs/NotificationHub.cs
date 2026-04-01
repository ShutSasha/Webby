using Microsoft.AspNetCore.SignalR;

namespace Webby.NotificationService.Hubs;

public class NotificationHub : Hub
{
   
   public override async Task OnConnectedAsync()
   {
      if (Context.User?.Identity?.IsAuthenticated != true)
      {
         await Clients.Caller.SendAsync("AuthError", "Unauthorized: Token is missing or expired");
         Context.Abort();
         return;
      }

      Console.WriteLine("CONNECTED: " + Context.UserIdentifier);
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