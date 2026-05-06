using AutoMapper;
using Microsoft.EntityFrameworkCore.Metadata.Internal;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Stream.Enums;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Helpers.Video;
using Webby.VideoService.Interfaces.Dto;
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
   private readonly IYouTubeSearchService _youtubeSearchService;
   public VideoService(IVideoRepository videoRepository, IStorageService storageService,
      ITagService tagService, IPlaylistService playlistService, UserGrpcService.UserGrpcServiceClient userClient, 
      IMapper mapper, IBackgroundTaskQueue queue,
      IServiceScopeFactory scopeFactory, IYouTubeSearchService youtubeSearchService)
   {
      _videoRepository = videoRepository;
      _storageService = storageService;
      _tagService = tagService;
      _playlistService = playlistService;
      _userClient = userClient;
      _mapper = mapper;
      _queue = queue;
      _scopeFactory = scopeFactory;
      _youtubeSearchService = youtubeSearchService;
   }

   public async Task<Video> GetVideoById(Guid videoId)
   {
      var video = await _videoRepository.FindById(videoId)
                  ?? throw new ApiException("Get video error", 404, "Video wasn't found");

      if (video.IsPrivate)
      {
         throw new ApiException("Get video error", 403, "Requested video is private");
      }

      if (video.VideoUploadStatus is 
          VideoStatus.Uploading or 
          VideoStatus.Canceled or 
          VideoStatus.Failed)
      {
         throw new ApiException("Get video error", 400, "Video is not uploaded yet");
      }

      if (!video.IsPublished)
      {
         throw new ApiException("Get video error", 400, "Video is not published yet");
      }

      return video;
   }

   public async Task<UploadVideoResponse> UploadVideoFile(Guid userId, UploadVideoRequest request)
   {
      if (request.VideoFile == null || request.VideoFile.Length == 0)
         throw new ArgumentException("File is empty");

      var videoId = Guid.NewGuid();

      var video = new Video
      {
         VideoId = videoId,
         Name = Path.GetFileNameWithoutExtension(request.VideoFile.FileName),
         UserId = userId,
         VideoUploadStatus =VideoStatus.Uploading,
         CreatedAt = DateTime.UtcNow,
         IsPrivate = false,
         IsPublished = false
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

      if (userId != video.UserId)
      {
         throw new ApiException("Publish video error", 403, "You can't publish this video");
      }
      
      await using var previewStream = request.PreviewFile.OpenReadStream();
      
      var previewUrl = await _storageService.UploadFileAsync(
         video.VideoId,
         "previews",
         request.PreviewFile.FileName,
         previewStream,
         request.PreviewFile.ContentType
      );

      video.Name = request.Name.Trim();
      video.Description = request.Description?.Trim();
      video.PreviewUrl = previewUrl;
      video.IsPrivate = request.IsPrivate;
      video.IsPublished = true;
      
      await _videoRepository.Update(video);
      
      if (request.VideoTags != null && request.VideoTags.Any())
      {
         await _tagService.EnsureCreateTags(request.VideoTags, video.VideoId);
      }
   }

   public async Task DeleteVideo(Guid userId, string videoId)
   {
      var (_, actualId) = ParseVideoPrefix(videoId);
      
      var video = await _videoRepository.FindById(Guid.Parse(actualId)) 
                  ?? throw new ApiException("Delete video error", 404, "Video wasn't found");

      if (video.UserId != userId)
      {
         throw new ApiException("Delete video error", 403, "You can't delete this video");
      }

      if (!string.IsNullOrEmpty(video.VideoUrl))
      {
         await _storageService.DeleteFileAsync(video.VideoUrl);
      }

      if (!string.IsNullOrEmpty(video.PreviewUrl))
      {
         await _storageService.DeleteFileAsync(video.PreviewUrl);
      }
      
      await _videoRepository.DeleteAsync(video.VideoId);
   }

   public async Task<VideoDto> GetVideoInformation(string videoId, Guid? userId)
   {
      var (platform, actualId) = ParseVideoPrefix(videoId);
      
      switch (platform)
      {
         case SearchVideoPlatforms.YouTube:
         {
            var youtubeVideoDto = await _youtubeSearchService.FindById(actualId);
            return youtubeVideoDto;
         }
         case SearchVideoPlatforms.Webby:
            break;
         default:
            throw new ApiException("Get video information error", 400, "incorrect platform type");
      }

      if (!Guid.TryParse(actualId, out var videoIdGuid))
      {
         throw new ApiException("Get video information error", 400, "The provided ID is not a valid GUID.");
      }
      
      var video = await _videoRepository.GetVideoInformationById(videoIdGuid)
                  ?? throw new ApiException("Get video information error", 404, "Video wasn't found");

      if (video.VideoUploadStatus is not VideoStatus.Ready)
      {
         throw new ApiException("Get video information error", 400, "Video is not uploaded yet");
      }

      var userResponse = await _userClient.GetUserByIdAsync(new GetUserRequest
      {
         UserId = video.UserId.ToString(),
         RequestUserId = userId?.ToString() ?? string.Empty
      });

      var videoTagsNames = await _tagService.GetTagNames(video.VideoTags?.ToList());

      return new VideoDto()
      {
         VideoId = PlatformPrefixesConstants.WebbyPrefix + videoIdGuid.ToString(),
         Name = video.Name,
         Views = video.Views,
         Description = video.Description,
         CreatedAt = video.CreatedAt,
         Duration = video.Duration,
         VideoUrl = video.VideoUrl,
         PreviewUrl = video.PreviewUrl,
         IsPrivate = video.IsPrivate,
         VideoTags = videoTagsNames,
         User = new UserVideoDto
         {
            UserId = userResponse.UserId,
            Username = userResponse.Username,
            AvatarUrl = userResponse.AvatarUrl,
            IsFollowed = userResponse.IsFollowed
         }
      };
   }

   public async Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId, Guid? requestedUserId, GetUserVideosRequest request)
   {
      var isOwner = userId == requestedUserId;

      var (videos, totalCount) = await _videoRepository
         .GetPaginatedUserVideos(userId, isOwner, request.Page, request.PageSize, request.ShouldShowDrafts);

      return new PagedResponse<VideoDto>
      {
         Items = MapToDto(videos),
         Page = request.Page,
         PageSize = request.PageSize,
         TotalCount = totalCount
      };
   }

   public async Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request)
   {
      var (_, actualId) = ParseVideoPrefix(request.VideoId);
      
      var video = await _videoRepository.GetVideoInformationById(Guid.Parse(actualId))
                  ?? throw new ApiException("Update video information error", 404, "Video wasn't found");

      if (video.UserId != userId)
         throw new ApiException("Update information error", 403, "You don't have permission for updating this video");
      
      video.Name = request.Name.Trim();
      video.Description = request.Description?.Trim();
      video.IsPrivate = request.IsPrivate;
      
      
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

      if (!video.IsPublished)
         video.IsPublished = true;

      await _videoRepository.Update(video);
   }

   public async Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchVideoOptions options)
   {
      options.SearchPlatform ??= SearchVideoPlatforms.Webby;
      
      if (options.SearchPlatform == SearchVideoPlatforms.YouTube)
      {
         var youtubeResponse = await _youtubeSearchService
            .SearchAsync(options.SearchText, options.PageSize, options.Page, options.NextPageToken);
         return youtubeResponse;
      }

      var skip = (options.Page - 1) * options.PageSize;
      
      List<Video> videos = [];
      var total = 0;
      var seed = options.ContentSeed;
      
      if (string.IsNullOrWhiteSpace(options.SearchText))
      {
         var (subscribedAuthorIds, historyTags) = await GetUserRecommendationContextAsync(requestUserId);
          
          (videos, total, seed) = await _videoRepository.GetRecommendedVideosAsync(
              currentVideoId: null,
              currentTags: [], 
              subscribedIds: subscribedAuthorIds,
              historyTags: historyTags,
              skip: skip,
              pageSize: options.PageSize,
              contentSeed: options.ContentSeed);
      }
      else
      {
          var (additionalConditional, parameters, predicate) = VideoQueryFilters.SearchVideoFilter(requestUserId);
          
          (videos, total) = await _videoRepository.SearchAsync(
             "Videos",
             "Name",
             options.SearchText,
             skip,
             options.PageSize,
             additionalConditional,
             parameters,
             predicateFactory: predicate
          );
      }
      
      if (videos.Count == 0)
      {
          return new PagedResponse<VideoDto>
          {
              Items = new List<VideoDto>(),
              TotalCount = total,
              Page = options.Page,
              PageSize = options.PageSize
          };
      }
      
      var items = await MapAndEnrichWithUsersAsync<VideoDto>(videos);
      
      return new PagedResponse<VideoDto>
      {
         Items = items,
         TotalCount = total,
         Page = options.Page,
         PageSize = options.PageSize,
         ContentSeed = seed
      };
   }

   public async Task<PagedResponse<VideoDto>> SearchVideoInPlaylist(Guid? requestUserId, Guid playlistId, SearchOptions searchOptions)
   {
       var searchText = searchOptions.SearchText?.Trim();
       bool hasSearch = !string.IsNullOrWhiteSpace(searchText);
       
       int dbSkip = hasSearch ? 0 : (searchOptions.Page - 1) * searchOptions.PageSize;
       int dbTake = hasSearch ? int.MaxValue : searchOptions.PageSize;
      
       var (items, total) = await _videoRepository.SearchVideosInPlaylistAsync(
           playlistId,
           requestUserId,
           searchText,
           dbSkip,
           dbTake
       );

       var resultItems = new List<VideoDto>();
       
       var youtubeIds = items
           .Where(pv => pv.VideoPlatform == VideoPlatform.YouTube && !string.IsNullOrEmpty(pv.ExternalVideoId))
           .Select(pv => pv.ExternalVideoId!)
           .ToList();

       var youtubeVideosDict = new Dictionary<string, VideoDto>();
       if (youtubeIds.Count != 0)
       {
           var ytList = await _youtubeSearchService.GetList(youtubeIds);
           youtubeVideosDict = ytList.ToDictionary(v => v.VideoId!);
       }
      
       foreach (var pv in items)
       {
           if (pv.VideoPlatform == VideoPlatform.Webby && pv.Video != null)
           {
               resultItems.Add(_mapper.Map<VideoDto>(pv.Video));
           }
           else if (pv.VideoPlatform == VideoPlatform.YouTube && !string.IsNullOrEmpty(pv.ExternalVideoId))
           {
               if (youtubeVideosDict.TryGetValue(PlatformPrefixesConstants.YouTubePrefix + pv.ExternalVideoId, out var ytVideo))
               {
                   if (hasSearch)
                   {
                       if (ytVideo.Name != null && ytVideo.Name.Contains(searchText!, StringComparison.OrdinalIgnoreCase))
                       {
                           resultItems.Add(ytVideo);
                       }
                   }
                   else
                   {
                       resultItems.Add(ytVideo);
                   }
               }
           }
       }
       
       if (hasSearch)
       {
           total = resultItems.Count;
           resultItems = resultItems
               .Skip((searchOptions.Page - 1) * searchOptions.PageSize)
               .Take(searchOptions.PageSize)
               .ToList();
       }

       return new PagedResponse<VideoDto>()
       {
           Items = resultItems,
           Page = searchOptions.Page,
           PageSize = searchOptions.PageSize,
           TotalCount = total
       };
   }

   public async Task<PagedResponse<PreviewVideoDto>> GetRecommendationVideos(string videoId, Guid? requestUserId,int contentSeed,
      int page = 1, int pageSize = 20)
   {
      VideoDto? watchingVideo = null;
      Guid? currentVideoGuid = null;


      var (platform, actualId) = ParseVideoPrefix(videoId);

      switch (platform)
      {
         case SearchVideoPlatforms.YouTube:
            watchingVideo = await _youtubeSearchService.FindById(videoId);
            break;
         case SearchVideoPlatforms.Webby:
            if (!Guid.TryParse(actualId, out var videoIdGuid))
            {
               throw new ApiException("Get recommendation videos error", 400, "Incorrect id format of local video");
            }
            
            watchingVideo = _mapper.Map<VideoDto>(await _videoRepository.GetVideoInformationById(videoIdGuid));
            break;
      }

      if (watchingVideo == null)
      {
         return new PagedResponse<PreviewVideoDto> { Items = [], Page = page, PageSize = pageSize, TotalCount = 0 };
      }

      var watchingVideoTags = watchingVideo.VideoTags ?? [];
      var skipVideos = (page - 1) * pageSize;
      
      var (subscribedAuthorIds, historyTags) = await GetUserRecommendationContextAsync(requestUserId);
      
      var (recommendedVideos, totalCount, seed) = await _videoRepository.GetRecommendedVideosAsync(
         currentVideoGuid,
         watchingVideoTags,
         subscribedAuthorIds,
         historyTags,
         skipVideos,
         pageSize,
         contentSeed);

      if (recommendedVideos.Count == 0)
      {
          return new PagedResponse<PreviewVideoDto>
          {
              Items = [],
              TotalCount = totalCount,
              Page = page,
              PageSize = pageSize
          };
      }
      
      var items = await MapAndEnrichWithUsersAsync<PreviewVideoDto>(recommendedVideos);

      return new PagedResponse<PreviewVideoDto>
      {
         Items = items,
         Page = page,
         PageSize = pageSize,
         TotalCount = totalCount,
         ContentSeed = seed
      };
   }

   public async Task<bool> CheckUploadStatus(Guid videoId)
   {
      var video = await _videoRepository.FindById(videoId);
      
      if (video == null) 
      {
         return false; 
      }
   
      return video.VideoUploadStatus switch
      {
         VideoStatus.Ready => true,
         VideoStatus.Canceled or VideoStatus.Failed or VideoStatus.Uploading => false,
         _ => false 
      };
   }

   public async Task IncrementVideoView(Guid requestUserId, string videoId)
   {

      var (platform, actualId) = ParseVideoPrefix(videoId);

      if (platform == SearchVideoPlatforms.YouTube)
         return;

      if (!Guid.TryParse(actualId, out var localVideoId))
         throw new ApiException("Increment video view error", 400, "Invalid local video id format type");
      
      if (await _videoRepository.FindUserView(requestUserId, localVideoId))
         return;

      var video = await _videoRepository.FindById(localVideoId);

      if (video == null)
         return;

      await _videoRepository.AddUserView(new UserView
      {
         UserId = requestUserId,
         VideoId = localVideoId
      });
      
      video.Views = await _videoRepository.CountUserView(localVideoId);
      await _videoRepository.Update(video);
   }

   public async Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId)
   {
      var privateVideos = await _videoRepository
         .GetByPredicate(v => v.IsPrivate 
                              && videoIds.Contains(v.VideoId) 
                              && v.UserId != requestUserId);

      return privateVideos?.Any() ?? false;
   }

   public async Task CancelVideoUploading(Guid requestUserId, Guid videoId)
   {
      var video = await _videoRepository.FindById(videoId)
                  ?? throw new ApiException("Cancel video uploading", 404, "Video wasn't found");

      if (requestUserId != video.UserId)
      {
         throw new ApiException("Publish video error", 403, "You can't publish this video");
      }
      
      switch (video.VideoUploadStatus)
      {
         case VideoStatus.Uploading:
            video.VideoUploadStatus = VideoStatus.Canceled;
            await _videoRepository.Update(video);
            break;

         case VideoStatus.Ready:
            if (!string.IsNullOrEmpty(video.VideoUrl))
               await _storageService.DeleteFileAsync(video.VideoUrl);

            await _videoRepository.DeleteAsync(videoId);
            break;

         case VideoStatus.Failed:
         case VideoStatus.Canceled:
            await _videoRepository.DeleteAsync(videoId);
            break;
      }
   }

   private List<VideoDto> MapToDto(IEnumerable<Video> videos) =>
      videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
   
   private async Task<(List<Guid> SubscribedIds, List<string> HistoryTags)> GetUserRecommendationContextAsync(Guid? requestUserId)
   {
      if (!requestUserId.HasValue)
         return ([], []);

      var subscriptionsResponse = await _userClient.GetUserSubscriptionIdsAsync(
         new GetUserSubscriptionIdsRequest { RequestUserId = requestUserId.ToString() });
            
      var subscribedIds = subscriptionsResponse.UserIds.Select(Guid.Parse).ToList();
      var historyTags = await _videoRepository.GetRecentUserViewTagsAsync(requestUserId.Value);

      return (subscribedIds, historyTags);
   }
   
   private async Task<List<TDto>> MapAndEnrichWithUsersAsync<TDto>(List<Video> videos) where TDto : IVideoDtoWithUser
   {
      if (videos.Count == 0) return [];

      var userIds = videos
         .Select(v => v.UserId.ToString())
         .Distinct()
         .ToList();

      var usersResponse = await _userClient.GetUsersByIdsAsync(new GetUsersRequest
      { 
         UserIds = { userIds } 
      });

      var usersDict = usersResponse.Users.ToDictionary(u => u.UserId, u => u);

      return videos.Select(v => 
      {
         var dto = _mapper.Map<TDto>(v);
        
         if (usersDict.TryGetValue(v.UserId.ToString(), out var userInfo))
         {
            dto.User = new UserVideoDto
            {
               UserId = userInfo.UserId,
               Username = userInfo.Username,
               AvatarUrl = userInfo.AvatarUrl,
            };
         }

         return dto;
      }).ToList();
   }

   private (SearchVideoPlatforms, string) ParseVideoPrefix(string prefixedId)
   {
      var parseResult = PlatformPrefixToPlatformConverter.ParseSearchVideoPlatform(prefixedId);

      if (parseResult == null)
      {
         throw new ApiException("Invalid ID format", 400, "Platform prefix is missing or unsupported");
      }

      var (platform, actualId) = parseResult.Value;

      return (platform, actualId);
      
   }
   
}