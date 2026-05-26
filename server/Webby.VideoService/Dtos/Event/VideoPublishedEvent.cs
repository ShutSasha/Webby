using Amazon.S3;
using Webby.VideoService.Dtos.Event.Enums;
using Webby.VideoService.Interfaces.Helpers;

namespace Webby.VideoService.Dtos.Event;

public record VideoPublishedEvent(Guid UserId, Guid VideoId) : IPlatformEvent
{
   public string EventType => "video_published";
   public EventOperation EventOperation => EventOperation.Absolute;
   public int Value { get; set; }
}