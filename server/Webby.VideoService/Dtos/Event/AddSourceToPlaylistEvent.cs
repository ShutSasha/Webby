using Webby.VideoService.Interfaces.Helpers;

namespace Webby.VideoService.Dtos.Event;

public record AddSourceToPlaylistEvent(Guid UserId,Guid playlistId) : IPlatformEvent
{
   public string EventType => "add_video_to_playlist";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}