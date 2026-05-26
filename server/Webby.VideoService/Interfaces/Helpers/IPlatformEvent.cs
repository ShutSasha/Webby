using Amazon.S3;
using Webby.VideoService.Dtos.Event.Enums;

namespace Webby.VideoService.Interfaces.Helpers;

public interface IPlatformEvent
{
   string EventType { get; }
   EventOperation EventOperation { get; }
   int Value { get; set; }
}