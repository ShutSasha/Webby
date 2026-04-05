using Microsoft.AspNetCore.SignalR;

namespace Webby.NotificationService.Providers;

public class UserProvider : IUserIdProvider
{
   public string GetUserId(HubConnectionContext connection)
   {
      return connection.User?.FindFirst("Id")?.Value;
   }
}