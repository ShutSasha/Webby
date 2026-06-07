using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Dtos.Event;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Stream;
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
using NewsStyleUriParser = System.NewsStyleUriParser;

namespace Webby.VideoService.Services;

public class VideoService : IVideoService
{
   private readonly IVideoRepository _videoRepository;
   private readonly IStorageService _storageService;
   private readonly ITagService _tagService;
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   private readonly IMapper _mapper;
   private readonly IBackgroundTaskQueue _queue;
   private readonly IServiceScopeFactory _scopeFactory;
   private readonly IYouTubeSearchService _youtubeSearchService;
   private readonly ITwitchSearchService _twitchSearchService;
   private readonly IEventPublisher _eventPublisher;
   private readonly ILogger<VideoService> _logger;
   
   public VideoService(IVideoRepository videoRepository, IStorageService storageService,
      ITagService tagService, UserGrpcService.UserGrpcServiceClient userClient, 
      IMapper mapper, IBackgroundTaskQueue queue,
      IServiceScopeFactory scopeFactory, IYouTubeSearchService youtubeSearchService,
      ITwitchSearchService twitchSearchService, IEventPublisher eventPublisher,
      ILogger<VideoService> logger)
   {
      _videoRepository = videoRepository;
      _storageService = storageService;
      _tagService = tagService;
      _userClient = userClient;
      _mapper = mapper;
      _queue = queue;
      _scopeFactory = scopeFactory;
      _youtubeSearchService = youtubeSearchService;
      _twitchSearchService = twitchSearchService;
      _eventPublisher = eventPublisher;
      _logger = logger;
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
      var userVideosCount = await _videoRepository.CountUserVideos(userId);

      if (userVideosCount >= LimitationConstants.FreeUploadVideosLimit)
      {
         var userPremiumStatus = await _userClient
            .GetPremiumStatusAsync(
               new GetPremiumStatusRequest()
               {
                  UserId = userId.ToString()
               });

         switch (userPremiumStatus.Status)
         {
            case PremiumStatus.Trial:
            case PremiumStatus.Active:
               break;
            case PremiumStatus.None:
               throw new ApiException("Upload file error", 403, $"You can't upload more than {LimitationConstants.FreeUploadVideosLimit} without premium status");
            case PremiumStatus.Expired:
               throw new ApiException("Upload file error",403,$"Your premium subscription have expired. You may upload only {LimitationConstants.FreeUploadVideosLimit} videos without premium status");
            default:
               throw new ApiException("Upload file error", 500, "Invalid premium status type");
         }
      }
      
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
      var (_, actualId) = ParseSystemPlatform(request.VideoId);
      
      var video = await _videoRepository.FindById(Guid.Parse(actualId))
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
      
      var platformEvent = new VideoPublishedEvent(userId, video.VideoId)
      {
         Value = await _videoRepository.CountUserVideos(userId)
      };
      
      await _eventPublisher.PublishAsync(platformEvent);
   }

   public async Task DeleteVideo(Guid userId, string videoId)
   {
      var (_, actualId) = ParseSystemPlatform(videoId);
      
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
      var (platform, actualId) = ParseSystemPlatform(videoId);
      
      switch (platform)
      {
         case SystemPlatforms.YouTube:
         {
            var youtubeVideoDto = await _youtubeSearchService.FindById(actualId);
            return youtubeVideoDto;
         }
         case SystemPlatforms.Webby:
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
      var videoDto = _mapper.Map<VideoDto>(video);
      
      videoDto.VideoTags = videoTagsNames;
      videoDto.User = new UserVideoDto()
      {
         UserId = userResponse.UserId,
         Username = userResponse.Username,
         AvatarUrl = userResponse.AvatarUrl,
         IsFollowed = userResponse.IsFollowed
      };
      
      return videoDto;
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
      var (_, actualId) = ParseSystemPlatform(request.VideoId);
      
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
      
      List<Video> videos;
      int total;
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
           .Where(pv => pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId))
           .Select(pv => pv.ExternalContentId!)
           .ToList();

       var twitchIds = items
           .Where(pv => pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId))
           .Select(pv => pv.ExternalContentId!)
           .ToList();

       var youtubeVideosDict = new Dictionary<string, VideoDto>();
       if (youtubeIds.Count != 0)
       {
           var ytList = await _youtubeSearchService.GetList(youtubeIds);
           youtubeVideosDict = ytList.ToDictionary(v => v.VideoId);
       }

       var twitchStreamsDict = new Dictionary<string, StreamDto>();
       if (twitchIds.Count != 0)
       {
           var twitchList = await _twitchSearchService.GetList(twitchIds);
           twitchStreamsDict = twitchList.ToDictionary(s => s.StreamerId);
       }
      
       foreach (var pv in items)
       {
           if (pv is { Platform: SystemPlatforms.Webby, Video: not null })
           {
               resultItems.Add(_mapper.Map<VideoDto>(pv.Video));
           }
           else if (pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId))
           {
               if (youtubeVideosDict.TryGetValue(PlatformPrefixesConstants.YouTubePrefix + pv.ExternalContentId, out var ytVideo))
               {
                   if (hasSearch)
                   {
                       if (ytVideo.Name.Contains(searchText!, StringComparison.OrdinalIgnoreCase))
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
           else if (pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId))
           {
               if (twitchStreamsDict.TryGetValue(PlatformPrefixesConstants.TwitchPrefix + pv.ExternalContentId, out var twitchStream))
               {
                  var streamAsVideo = _mapper.Map<VideoDto>(twitchStream);
                  streamAsVideo.CreatedAt = pv.CreatedAt;
                   if (hasSearch)
                   {
                       if (streamAsVideo.Name.Contains(searchText!, StringComparison.OrdinalIgnoreCase))
                       {
                           resultItems.Add(streamAsVideo);
                       }
                   }
                   else
                   {
                       resultItems.Add(streamAsVideo);
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


      var (platform, actualId) = ParseSystemPlatform(videoId);

      switch (platform)
      {
         case SystemPlatforms.YouTube:
            watchingVideo = await _youtubeSearchService.FindById(actualId);
            break;
         case SystemPlatforms.Webby:
            if (!Guid.TryParse(actualId, out var videoIdGuid))
            {
               throw new ApiException("Get recommendation videos error", 400, "Incorrect id format of local video");
            }
            
            watchingVideo = _mapper.Map<VideoDto>(await _videoRepository.GetVideoInformationById(videoIdGuid));
            break;
         case SystemPlatforms.Twitch:
         {
            watchingVideo = _mapper.Map<VideoDto>(await _twitchSearchService.FindById(actualId));
            break;
         }
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

   public async Task<bool> CheckUploadStatus(string videoId)
   {
      var (_, actualId) = ParseSystemPlatform(videoId);
      
      var video = await _videoRepository.FindById(Guid.Parse(actualId));
      
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

      var (platform, actualId) = ParseSystemPlatform(videoId);

      if (platform == SystemPlatforms.YouTube)
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

   //TODO: Refactor structuring
   public async Task<(List<VideoDto>, List<string>)> GetVideoRange(List<string> ids)
   {
      
      if (ids.Count == 0)
      {
         return ([],[]);
      }

      var webbyIds = new List<Guid>();
      var youtubeIds = new List<string>();
      var twitchIds = new List<string>();
      var unavailableVideos = new List<string>();
      
      foreach (var id in ids)
      {
         var parseResult = PlatformPrefixToPlatformConverter.ParseSystemPlatform(id);
         
         if (parseResult == null)
            continue;

         var (platform, actualId) = parseResult.Value;
         

         switch (platform)
         {
            case SystemPlatforms.Webby:
               if (Guid.TryParse(actualId, out var guidId))
               {
                  webbyIds.Add(guidId);
               }
               break;
            case SystemPlatforms.YouTube:
               youtubeIds.Add(actualId);
               break;
            case SystemPlatforms.Twitch:
               twitchIds.Add(actualId);
               break;
         }
      }
      
      var fetchedVideosDict = new Dictionary<string, VideoDto>(StringComparer.OrdinalIgnoreCase);
      
      if (webbyIds.Count != 0)
      {
         var webbyVideos = await _videoRepository
            .GetByPredicate(v => webbyIds.Contains(v.VideoId) && 
                                 !v.IsPrivate &&
                                 v.IsPublished);
         
         if (webbyVideos != null && webbyVideos.Any())
         {
            var enrichedWebbyVideos = await MapAndEnrichWithUsersAsync<VideoDto>(webbyVideos.ToList());
            foreach (var video in enrichedWebbyVideos)
            {
               fetchedVideosDict[video.VideoId] = video; 
            }
         }
      }
      
      if (youtubeIds.Count != 0)
      {
         var youtubeVideos = await _youtubeSearchService.GetList(youtubeIds);
         
         if (youtubeVideos != null)
         {
            foreach (var video in youtubeVideos)
            {
               fetchedVideosDict[video.VideoId] = video;
            }
         }
      }
      
      if (twitchIds.Count != 0)
      {
         var twitchVideos = await _twitchSearchService.GetList(twitchIds); 
         
         if (twitchVideos != null)
         {
            foreach (var stream in twitchVideos)
            {
               var video = _mapper.Map<VideoDto>(stream);
               fetchedVideosDict[video.VideoId] = video;
            }
         }
      }

      var resultItems = new List<VideoDto>();
      
      foreach (var id in ids)
      {
         if (fetchedVideosDict.TryGetValue(id, out var videoDto))
         {
            resultItems.Add(videoDto);
         }
         else
         {
            unavailableVideos.Add(id);
         }
      }

      return (resultItems,unavailableVideos);
   }

   public async Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId)
   {
      var privateVideos = await _videoRepository
         .GetByPredicate(v => v.IsPrivate 
                              && videoIds.Contains(v.VideoId) 
                              && v.UserId != requestUserId);

      return privateVideos?.Any() ?? false;
   }

   public async Task CancelVideoUploading(Guid requestUserId, string videoId)
   {
      var (_, actualId) = ParseSystemPlatform(videoId);
      
      var video = await _videoRepository.FindById(Guid.Parse(actualId))
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

            await _videoRepository.DeleteAsync(video.VideoId);
            break;

         case VideoStatus.Failed:
         case VideoStatus.Canceled:
            await _videoRepository.DeleteAsync(video.VideoId);
            break;
      }
   }

   private List<VideoDto> MapToDto(IEnumerable<Video> videos) =>
      videos.Select(v => _mapper.Map<VideoDto>(v)).ToList();
   
   private async Task<(List<Guid> SubscribedIds, List<string> HistoryTags)> GetUserRecommendationContextAsync(Guid? requestUserId)
   {
      if (!requestUserId.HasValue)
         return ([], []);

      var subscribedIds = new List<Guid>();

      try
      {
         var response = await _userClient.GetUserSubscriptionIdsAsync(new GetUserSubscriptionIdsRequest
            { RequestUserId = requestUserId.ToString() });

         if (response.UserIds != null)
         {
            subscribedIds = response.UserIds
               .Where(id => Guid.TryParse(id, out _))
               .Select(Guid.Parse)
               .ToList();
         }
      }
      catch (RpcException ex)
      {
         _logger.LogWarning(ex, "Failed to fetch subscriptions for user {UserId}", requestUserId);
      }

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

   private (SystemPlatforms, string) ParseSystemPlatform(string prefixedId)
   {
      var parseResult = PlatformPrefixToPlatformConverter.ParseSystemPlatform(prefixedId);

      if (parseResult == null)
      {
         throw new ApiException("Invalid ID format", 400, "Platform prefix is missing or unsupported");
      }

      var (platform, actualId) = parseResult.Value;

      return (platform, actualId);
      
   }
   
}