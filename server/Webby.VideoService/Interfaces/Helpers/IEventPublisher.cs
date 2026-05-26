namespace Webby.VideoService.Interfaces.Helpers;

public interface IEventPublisher
{
   Task PublishAsync<T>(T @event) where T : IPlatformEvent;
}