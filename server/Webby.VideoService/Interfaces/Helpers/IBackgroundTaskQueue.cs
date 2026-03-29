namespace Webby.VideoService.Interfaces.Helpers;

public interface IBackgroundTaskQueue
{
   void Enqueue(Func<CancellationToken, Task> workItem);

   Task<Func<CancellationToken, Task>> DequeueAsync(CancellationToken cancellationToken);
}