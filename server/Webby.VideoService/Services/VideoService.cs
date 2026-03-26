using AutoMapper;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Helpers.Video;
using Webby.VideoService.Interfaces.Helpers;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;
using Webby.VideoService.Services.Background;

namespace Webby.VideoService.Services;

public class VideoService : IVideoService
{
   private readonly IVideoRepository _videoRepository;
   private readonly IStorageService _storageService;
   private readonly ITagService _tagService;
   private readonly IPlaylistService _playlistService;
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   private readonly IMapper _mapper;
   private readonly IBackgroundTaskQueue _queue;
   private readonly IServiceScopeFactory _scopeFactory;
   public VideoService(IVideoRepository videoRepository, IStorageService storageService, ITagService tagService, IPlaylistService playlistService, UserGrpcService.UserGrpcServiceClient userClient, IMapper mapper, IBackgroundTaskQueue queue, IServiceScopeFactory scopeFactory)
   {
      _videoRepository = videoRepository;
      _storageService = storageService;
      _tagService = tagService;
      _playlistService = playlistService;
      _userClient = userClient;
      _mapper = mapper;
      _queue = queue;
      _scopeFactory = scopeFactory;
   }

   public async Task<UploadVideoResponse> UploadVideoFile(Guid userId, UploadVideoRequest request)
   {
      if (request.VideoFile == null || request.VideoFile.Length == 0)
         throw new ArgumentException("File is empty");

      var videoId = Guid.NewGuid();

      var video = new Video
      {
         VideoId = videoId,
         Name = request.VideoFile.FileName,
         UserId = userId,
         VideoUploadStatus =VideoStatus.Uploading,
         CreatedAt = DateTime.UtcNow,
         IsPrivate = false
      };

      await _videoRepository.Add(video);
      
      var tempPath = Path.Combine(Path.GetTempPath(), $"{videoId}_{request.VideoFile.FileName}");

      await using (var stream = File.Create(tempPath))
      {
         await request.VideoFile.CopyToAsync(stream);
      }
      
      _queue.Enqueue(async token =>
      {
         using var scope = _scopeFactory.CreateScope();

         var processor = scope.ServiceProvider
            .GetRequiredService<VideoUploadProcessor>();

         await processor.ProcessUpload(
            videoId,
            tempPath,
            request.VideoFile.ContentType,
            token);
      });

      return _mapper.Map<UploadVideoResponse>(video);
   }

   public async Task CreateVideo(Guid userId, CreateVideoRequest request)
   {
      var video = await _videoRepository.FindById(request.VideoId)
                  ?? throw new ApiException("Create video error", 404, "Video wasn't found");
      
      await using var previewStream = request.PreviewFile.OpenReadStream();
      
      var previewUrl = await _storageService.UploadFileAsync(
         video.VideoId,
         "previews",
         request.PreviewFile.FileName,
         previewStream,
         request.PreviewFile.ContentType
      );

      video.Name = request.Name;
      video.Description = request.Description;
      video.PreviewUrl = previewUrl;
      video.IsPrivate = request.IsPrivate;
      
      await _videoRepository.Update(video);
      
      if (request.VideoTags != null && request.VideoTags.Any())
      {
         await _tagService.EnsureCreateTags(request.VideoTags, video.VideoId);
      }
      
      if (request.PlaylistId.HasValue)
      {
         await _playlistService.AttachVideoToPlaylist(
            request.PlaylistId.Value,
            [video.VideoId],
            userId
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
         Description = video.Description,
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
            [request.VideoId],
            userId);
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

   public async Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchOptions options)
   {
      var skip = (options.Page - 1) * options.PageSize;
      var (additionalConditional, parameters, predicate) = VideoQueryFilters.SearchVideoFilter(requestUserId);
      
      var (videos, total) = await _videoRepository.SearchAsync(
         "Videos",
         "Name",
         options.SearchText,
         skip,
         options.PageSize,
         additionalConditional,
         parameters,
         predicateFactory: predicate
         );

      var items = videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
      
      return new PagedResponse<VideoDto>
      {
         Items = items,
         TotalCount = total,
         Page = options.Page,
         PageSize = options.PageSize
      };
   }

   public async Task<PagedResponse<VideoDto>> SearchVideoInPlaylist(Guid? requestUserId,Guid playlistId, SearchOptions searchOptions)
   {
      var skip = (searchOptions.Page - 1) * searchOptions.PageSize;
      
      var (items, total) = await _videoRepository.SearchVideosInPlaylistAsync(
         playlistId,
         requestUserId,
         searchOptions.SearchText,
         skip,
         searchOptions.PageSize
         );

      return new PagedResponse<VideoDto>()
      {
         Items = items.Select(v => _mapper.Map<VideoDto>(v)).ToList(),
         Page = searchOptions.Page,
         PageSize = searchOptions.PageSize,
         TotalCount = total
      };
   }

   public async Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId)
   {
      var privateVideos = await _videoRepository
         .GetByPredicate(v => v.IsPrivate 
                              && videoIds.Contains(v.VideoId) 
                              && v.UserId != requestUserId);

      return privateVideos?.Any() ?? false;
   }

   private List<VideoDto> MapToDto(IEnumerable<Video> videos) =>
      videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
   
}