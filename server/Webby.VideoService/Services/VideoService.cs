using AutoMapper;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
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
   private readonly YoutubeSearchService _youtubeSearchService;
   private readonly TwitchSearchService _twitchSearchService;
   public VideoService(IVideoRepository videoRepository, IStorageService storageService, ITagService tagService, IPlaylistService playlistService, UserGrpcService.UserGrpcServiceClient userClient, IMapper mapper, IBackgroundTaskQueue queue, IServiceScopeFactory scopeFactory, YoutubeSearchService youtubeSearchService, TwitchSearchService twitchSearchService)
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
      _twitchSearchService = twitchSearchService;
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

      video.Name = request.Name;
      video.Description = request.Description;
      video.PreviewUrl = previewUrl;
      video.IsPrivate = request.IsPrivate;
      video.IsPublished = true;
      
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

      return new GetVideoInformationResponse
      {
         VideoId = videoId.ToString(),
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
            UserId = userResponse.UserId,
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

   public async Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchVideoOptions options)
{
   if (options.SearchPlatform == SearchPlatforms.YouTube)
   {
      var youtubeResponse = await _youtubeSearchService.SearchAsync(options);
      return youtubeResponse;
   }

   if (options.SearchPlatform == SearchPlatforms.Twitch)
   {
      var twitchResponse = await _twitchSearchService.SearchAsync(options);
      return twitchResponse;
   }
   
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
   
   if (!videos.Any())
   {
       return new PagedResponse<VideoDto>
       {
           Items = new List<VideoDto>(),
           TotalCount = total,
           Page = options.Page,
           PageSize = options.PageSize
       };
   }
   
   var userIds = videos
      .Select(v => v.UserId.ToString())
      .Distinct()
      .ToList();
   
   var users = await _userClient.GetUsersByIdsAsync(new GetUsersRequest
   { 
       UserIds = { userIds } 
   });

   var usersDict = users.Users.ToDictionary(u => u.UserId, u => u);
   
   var items = videos.Select(v => 
   {
       var dto = _mapper.Map<VideoDto>(v);
       
       var currentUserIdStr = v.UserId.ToString();
       if (usersDict.TryGetValue(currentUserIdStr, out var userInfo))
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

   
   //TODO: refactor global search method
   public async Task<GlobalSearchVideoResponse> GlobalSearchVideos(Guid? requestUserId, GlobalSearchOptions searchOptions)
   {
      if (searchOptions.SectionType is not (SearchSections.Playlists or SearchSections.Videos))
      {
         throw new ApiException("Global video search error",400,"Incorrect search section type");
      }

      if (searchOptions.SectionType == SearchSections.Streams)
      {
         var twitchStreams = await _twitchSearchService.SearchAsync(new SearchVideoOptions
         {
            Page = searchOptions.Page,
            SearchPlatform = SearchPlatforms.Twitch,
            SearchText = searchOptions.SearchText,
            PageSize = 5
         });

         return new GlobalSearchVideoResponse()
         {
            TwitchStreams = new SearchSection<VideoDto>()
            {
               Items = twitchStreams.Items,
               NextPageToken = twitchStreams.NextPageToken,
               Page = searchOptions.Page
            }
         };
      }
      
      var skip = (searchOptions.Page - 1) * 5;
      var (additionalConditional, parameters, predicate) = VideoQueryFilters.SearchVideoFilter(requestUserId);

      var webbySearchTask = _videoRepository.SearchAsync(
         "Videos",
         "Name",
         searchOptions.SearchText,
         skip,
         searchOptions.PageSize,
         additionalConditional,
         parameters,
         predicateFactory: predicate
      );

      var youtubeSearchTask = _youtubeSearchService.SearchAsync(new SearchVideoOptions
      {
         Page = searchOptions.Page,
         PageSize = 5,
         SearchPlatform = SearchPlatforms.YouTube,
         SearchText = searchOptions.SearchText
      });

      await Task.WhenAll(webbySearchTask, youtubeSearchTask);

      var (webbyVideos, total) = await webbySearchTask;
      var youtubeVideos = await youtubeSearchTask;

      return new GlobalSearchVideoResponse
      {
         WebbyVideos = new SearchSection<VideoDto>
         {
            Items = MapToDto(webbyVideos),
            Page = searchOptions.Page,
            TotalCount = total,
         },
         YouTubeVideos = new SearchSection<VideoDto>
         {
            Items = youtubeVideos.Items,
            NextPageToken = youtubeVideos.NextPageToken,
            Page = searchOptions.Page,
            TotalCount = youtubeVideos.TotalCount
         }
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
   
}