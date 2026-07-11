using System.Diagnostics;
using Google.Apis.Auth.OAuth2.Flows;
using Google.Cloud.Vision.V1;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services;

public class VideoModerationService(
    ImageAnnotatorClient visionClient,
    ILogger<VideoModerationService> logger)
    : IVideoModerationService
{
    private readonly ILogger<VideoModerationService> _logger = logger;

   public async Task<bool> IsVideoSafeAsync(string videoFilePath, int durationSeconds, CancellationToken token)
    {
        var timestamps = new[] { durationSeconds * 0.1, durationSeconds * 0.5, durationSeconds * 0.9 };

        foreach (var seconds in timestamps)
        {
            if (token.IsCancellationRequested) break;

            var framePath = Path.Combine(Path.GetTempPath(), $"{Guid.NewGuid()}.jpg");

            try
            {
                await ExtractFrameAsync(videoFilePath, framePath, seconds, token);

                if (!File.Exists(framePath)) continue;

                var image = await Image.FromFileAsync(framePath);
                var safeSearch = await visionClient.DetectSafeSearchAsync(image);

                if (safeSearch.Adult is Likelihood.Likely or Likelihood.VeryLikely ||
                    safeSearch.Violence is Likelihood.Likely or Likelihood.VeryLikely)
                {
                    return false;
                }
            }
            catch (Exception ex)
            {
                _logger.LogError("Moderation checking error: {Message}",ex.Message);
            }
            finally
            {
                if (File.Exists(framePath))
                {
                    File.Delete(framePath);
                }
            }
        }

        return true;
    }

    private async Task ExtractFrameAsync(string videoPath, string outputPath, double timeSeconds, CancellationToken token)
    {
        var processInfo = new ProcessStartInfo
        {
            FileName = "ffmpeg",
            Arguments = $"-y -v error -ss {timeSeconds.ToString("0.##", System.Globalization.CultureInfo.InvariantCulture)} -i \"{videoPath}\" -vframes 1 -q:v 2 \"{outputPath}\"",
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true
        };

        using var process = new Process { StartInfo = processInfo };
        process.Start();

        var outputTask = process.StandardOutput.ReadToEndAsync();
        var errorTask = process.StandardError.ReadToEndAsync();

        await Task.WhenAll(process.WaitForExitAsync(token), outputTask, errorTask);

        if (process.ExitCode != 0)
        {
            var errorOutput = await errorTask;
            throw new InvalidOperationException($"FFmpeg frame extraction failed. Exit code: {process.ExitCode}. Output: {errorOutput}");
        }
    }
}