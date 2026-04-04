using NReco.VideoInfo;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services.Background;

public class VideoUploadProcessor(IVideoRepository videoRepository, IStorageService storageService, ILogger<VideoUploadProcessor> logger)
{
   public async Task ProcessUpload(
      Guid videoId,
      string filePath,
      string contentType,
      CancellationToken token)
   {
      var video = await videoRepository.FindById(videoId);
      string? videoFileUrl = null;
      bool isVideoDeleted = false;
   
      try
      {
         var ffProbe = new FFProbe();
         var videoInfo = ffProbe.GetMediaInfo(filePath);
         video!.Duration = (long)Math.Round(videoInfo.Duration.TotalSeconds);

         logger.LogInformation($"Uploading file to storage");
         await using var stream = File.OpenRead(filePath);

         videoFileUrl = await storageService.UploadFileAsync(
            videoId,
            "videos",
            Path.GetFileName(filePath),
            stream,
            contentType
         );

         video.VideoUrl = videoFileUrl;
         await videoRepository.Update(video);

         File.Delete(filePath);
      }
      catch (Exception ex)
      {
         if (!string.IsNullOrEmpty(videoFileUrl))
         {
            try
            {
               await storageService.DeleteFileAsync(videoFileUrl);
               isVideoDeleted = true;
            }
            catch (Exception deleteEx)
            {
               throw new ApiException("File rollback delete error", 500, deleteEx.Message);
            }
         }

         video!.VideoUploadStatus = VideoStatus.Failed;
         await videoRepository.Update(video);
      }
      finally
      {
         if (!string.IsNullOrEmpty(videoFileUrl))
         {
            await videoRepository.ReloadAsync(video!);
         
            logger.LogInformation($"Check video with status {video!.VideoUploadStatus}");
            if (video is { VideoUploadStatus: VideoStatus.Canceled })
            {
               logger.LogInformation($"Catch canceled video");
               if (isVideoDeleted)
               {
                  await storageService.DeleteFileAsync(videoFileUrl); 
               }
               await videoRepository.DeleteAsync(videoId);
            }
            else
            {
               video.VideoUploadStatus = VideoStatus.Ready;
               await videoRepository.Update(video);
            }
         }
      }
   }
}