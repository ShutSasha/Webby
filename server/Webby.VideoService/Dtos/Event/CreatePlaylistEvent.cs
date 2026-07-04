using Webby.VideoService.Interfaces.Helpers;

namespace Webby.VideoService.Dtos.Event;

public record CreatePlaylistEvent(Guid UserId, Guid playlistId): IPlatformEvent
{
   public string EventType => "create_playlist";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}