using Webby.UserService.Interfaces.Helpers;

namespace Webby.UserService.Dtos.Event;

public record GetModeratorRoleEvent(Guid UserId) : IPlatformEvent
{
   public string EventType => "get_role_moderator";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}