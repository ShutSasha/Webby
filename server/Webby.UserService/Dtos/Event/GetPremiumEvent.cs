using Webby.UserService.Interfaces.Helpers;

namespace Webby.UserService.Dtos.Event;

public record GetPremiumEvent(Guid UserId) : IPlatformEvent
{
   public string EventType => "get_premium";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}