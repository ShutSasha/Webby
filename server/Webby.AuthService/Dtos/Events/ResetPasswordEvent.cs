using Webby.AuthService.Interfaces.Helpers;

namespace Webby.AuthService.Dtos.Events;

public record ResetPasswordEvent(Guid UserId) : IPlatformEvent
{
   public string EventType => "reset_password";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}