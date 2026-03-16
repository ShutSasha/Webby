using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;

namespace Webby.VideoService.Services;

public class VideoService : IVideoService
{
   private readonly IVideoRepository _videoRepository;
   private readonly IStorageService _storageService;
   private readonly ITagService _tagService;
   private readonly IPlaylistService _playlistService;
   public VideoService(IVideoRepository videoRepository, IStorageService storageService, ITagService tagService, IPlaylistService playlistService)
   {
      _videoRepository = videoRepository;
      _storageService = storageService;
      _tagService = tagService;
      _playlistService = playlistService;
   }
   
   public async Task CreateVideo(Guid userId, CreateVideoRequest request)
   {
      var videoId = Guid.NewGuid();
      
      string videoUrl;
      await using (var videoStream = request.VideoFile.OpenReadStream())
      {
         videoUrl = await _storageService.UploadFileAsync(
            videoId,
            "videos",
            request.VideoFile.FileName,
            videoStream,
            request.VideoFile.ContentType
         );
      }

      await using var previewStream = request.PreviewFile.OpenReadStream();
      var previewUrl = await _storageService.UploadFileAsync(
         videoId,
         "previews",
         request.PreviewFile.FileName,
         previewStream,
         request.PreviewFile.ContentType
      );
      
      var video = new Video
      {
         VideoId = videoId,
         UserId = userId,
         Name = request.Name,
         Description = request.Description,
         Views = 0,
         CreatedAt = DateTime.UtcNow,
         VideoUrl = videoUrl,
         PreviewUrl = previewUrl,
         IsPrivate = request.IsPrivate
      };

      await _videoRepository.Add(video);
      
      if (request.VideoTags != null && request.VideoTags.Any())
      {
         await _tagService.EnsureCreateTags(request.VideoTags, videoId);
      }
      
      if (request.PlaylistId.HasValue)
      {
         await _playlistService.AttachVideoToPlaylist(
            request.PlaylistId.Value,
            [videoId]
         );
      }
      
   }

   public async Task DeleteVideo(Guid userId, Guid videoId)
   {
      var video = await _videoRepository.FindById(videoId);

      if (video == null)
      {
         throw new ApiException("Delete video error", 404, "Video wasn't found");
      }

      if (video.UserId != userId)
      {
         throw new ApiException("Delete video error", 403, "You can't delete this video");
      }
      
      await _storageService.DeleteFileAsync(video.VideoUrl);
      await _storageService.DeleteFileAsync(video.PreviewUrl);
      await _videoRepository.DeleteAsync(videoId);
   }
}