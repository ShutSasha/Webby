using Webby.VideoService.Interfaces.Helpers;

namespace Webby.VideoService.Dtos.Event;

public record VideoPublishedEvent(Guid UserId, Guid VideoId) : IPlatformEvent
{
   public string EventType => "video_published";
   public bool IsIncrementOperation => false;
   public int Value { get; set; }
}