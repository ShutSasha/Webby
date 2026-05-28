namespace Webby.VideoService.Interfaces.Helpers;

public interface IPlatformEvent
{
   string EventType { get; }
   bool IsIncrementOperation { get; }
   int Value { get; set; }
}