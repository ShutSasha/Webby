using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services.Background;

public class VideoUploadProcessor(IVideoRepository videoRepository, IStorageService storageService)
{
   public async Task ProcessUpload(
      Guid videoId,
      string filePath,
      string contentType,
      CancellationToken token)
   {
      var video = await videoRepository.FindById(videoId);

      try
      {
         await using var stream = File.OpenRead(filePath);

         var url = await storageService.UploadFileAsync(
            videoId,
            "videos",
            Path.GetFileName(filePath),
            stream,
            contentType
         );

         video.VideoUrl = url;
         video.VideoUploadStatus = VideoStatus.Ready;

         await videoRepository.Update(video);

         File.Delete(filePath);
      }
      catch (Exception)
      {
         video.VideoUploadStatus = VideoStatus.Failed;
         await videoRepository.Update(video);
      }
   }
}