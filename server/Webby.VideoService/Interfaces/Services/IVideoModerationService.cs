namespace Webby.VideoService.Interfaces.Services;

public interface IVideoModerationService
{
   Task<bool> IsVideoSafeAsync(string videoFilePath, int durationSeconds, CancellationToken token);
}