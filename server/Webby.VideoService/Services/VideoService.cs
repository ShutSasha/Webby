using AutoMapper;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Response;
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
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   private readonly IMapper _mapper;
   public VideoService(IVideoRepository videoRepository, IStorageService storageService, ITagService tagService, IPlaylistService playlistService, UserGrpcService.UserGrpcServiceClient userClient, IMapper mapper)
   {
      _videoRepository = videoRepository;
      _storageService = storageService;
      _tagService = tagService;
      _playlistService = playlistService;
      _userClient = userClient;
      _mapper = mapper;
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
      var video = await _videoRepository.FindById(videoId) 
                  ?? throw new ApiException("Delete video error", 404, "Video wasn't found");

      if (video.UserId != userId)
      {
         throw new ApiException("Delete video error", 403, "You can't delete this video");
      }
      
      await _storageService.DeleteFileAsync(video.VideoUrl);
      await _storageService.DeleteFileAsync(video.PreviewUrl);
      await _videoRepository.DeleteAsync(videoId);
   }
   
   public async Task<GetVideoInformationResponse> GetVideoInformation(Guid videoId, Guid? userId)
   {
      var video = await _videoRepository.GetVideoInformationById(videoId)
                  ?? throw new ApiException("Get video information error", 404, "Video wasn't found");

      var userResponse = await _userClient.GetUserByIdAsync(new GetUserRequest
      {
         UserId = video.UserId.ToString(),
         RequestUserId = userId?.ToString() ?? string.Empty
      });

      var videoTagsNames = await _tagService.GetTagNames(video.VideoTags?.ToList());

      return new GetVideoInformationResponse
      {
         VideoId = videoId,
         Name = video.Name,
         Views = video.Views,
         CreatedAt = video.CreatedAt,
         VideoUrl = video.VideoUrl,
         PreviewUrl = video.PreviewUrl,
         IsPrivate = video.IsPrivate,
         VideoTags = videoTagsNames,
         User = new UserVideoDto
         {
            UserId = Guid.Parse(userResponse.UserId),
            Username = userResponse.Username,
            AvatarUrl = userResponse.AvatarUrl,
            IsFollowed = userResponse.IsFollowed
         }
      };
   }

   public async Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId, Guid? requestUserId, int page, int pageSize)
   {
      page = page <= 0 ? 1 : page;
      pageSize = pageSize <= 0 ? 10 : pageSize;

      var isOwner = userId == requestUserId;

      var (videos, totalCount) = await _videoRepository
         .GetPaginatedUserVideos(userId, isOwner, page, pageSize);

      return new PagedResponse<VideoDto>
      {
         Items = MapToDto(videos),
         Page = page,
         PageSize = pageSize,
         TotalCount = totalCount
      };
   }

   public async Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request)
   {
      var video = await _videoRepository.GetVideoInformationById(request.VideoId)
                  ?? throw new ApiException("Update video information error", 404, "Video wasn't found");

      if (video.UserId != userId)
         throw new ApiException("Update information error", 403, "You don't have permission for updating this video");
      
      
      video.Name = request.Name;
      video.Description = request.Description;
      video.IsPrivate = request.IsPrivate;
      
      if (request.PlaylistId != null)
      {
         await _playlistService.AttachVideoToPlaylist(
            request.PlaylistId.Value,
            new List<Guid> { request.VideoId });
      }
      
      if (request.VideoFile != null)
      {
         if (!string.IsNullOrEmpty(video.VideoUrl))
         {
            await _storageService.DeleteFileAsync(video.VideoUrl);
         }

         await using var stream = request.VideoFile.OpenReadStream();

         var videoUrl = await _storageService.UploadFileAsync(
            video.VideoId,
            "videos",
            request.VideoFile.FileName,
            stream,
            request.VideoFile.ContentType);

         video.VideoUrl = videoUrl;
      }
      
      if (request.PreviewFile != null)
      {
         if (!string.IsNullOrEmpty(video.PreviewUrl))
         {
            await _storageService.DeleteFileAsync(video.PreviewUrl);
         }

         await using var stream = request.PreviewFile.OpenReadStream();

         var previewUrl = await _storageService.UploadFileAsync(
            video.VideoId,
            "previews",
            request.PreviewFile.FileName,
            stream,
            request.PreviewFile.ContentType);

         video.PreviewUrl = previewUrl;
      }
      
      if (request.VideoTags != null)
      {
         await _tagService.SyncVideoTags(video, request.VideoTags);
      }

      await _videoRepository.Update(video);
   }

   public async Task<PagedResponse<VideoDto>> SearchVideo(SearchVideoOptions options)
   {
      var skip = (options.Page - 1) * options.PageSize;

      var (videos, total) = await _videoRepository.SearchAsync(
         options.SearchText,
         skip,
         options.PageSize);

      var items = videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
      

      return new PagedResponse<VideoDto>
      {
         Items = items,
         TotalCount = total,
         Page = options.Page,
         PageSize = options.PageSize
      };
   }


   private List<VideoDto> MapToDto(IEnumerable<Video> videos) =>
      videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
}