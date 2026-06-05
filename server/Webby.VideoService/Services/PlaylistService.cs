using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Playlist;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services;

public class PlaylistService : IPlaylistService
{
   private readonly IPlaylistRepository _playlistRepository;
   private readonly IMapper _mapper;
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   private readonly IVideoRepository _videoRepository;
   private readonly IYouTubeSearchService _youtubeSearchService;
   private readonly ITwitchSearchService _twitchSearchService;

   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper,
      UserGrpcService.UserGrpcServiceClient userClient, IVideoRepository videoRepository,
      IYouTubeSearchService youtubeSearchService, ITwitchSearchService twitchSearchService)
   {
      _playlistRepository = playlistRepository;
      _mapper = mapper;
      _userClient = userClient;
      _videoRepository = videoRepository;
      _youtubeSearchService = youtubeSearchService;
      _twitchSearchService = twitchSearchService;
   }

   public async Task<Playlist> GetPlaylistById(Guid playlistId)
   {
      var playlist = await _playlistRepository.FindById(playlistId)
                     ?? throw new ApiException("Get playlist error", 404, "Playlist wasn't found");
      
      if (playlist.IsPrivate)
         throw new ApiException("Get playlist error",403,"Playlist is private");

      return playlist;
   }

   public async Task<PlaylistDto> CreatePlaylist(Guid userId, CreatePlaylistRequest request)
   {
      var playlist = new Playlist()
      {
         PlaylistId = Guid.NewGuid(),
         IsPrivate = request.IsPrivate,
         CreatedAt = DateTime.UtcNow,
         Name = request.Name,
         UserId = userId
      };

      await _playlistRepository.Add(playlist);

      return _mapper.Map<PlaylistDto>(playlist);
   }

   public async Task<PagedResponse<PlaylistPreviewDto>> GetUserPlaylists(Guid? requestUserId, string videoId, Guid userId, GetUserPlaylistsRequest request)
   {
      string actualId = default;
      
      var skip = (request.Page - 1) * request.PageSize;

      var (additionalCondition, parameters, predicate)
         = PlaylistSearchFilter.SearchUserPlaylistsFilter(userId, requestUserId == userId,request.ShouldShowEmptyPlaylists);

      var playlists = await _playlistRepository.SearchAsync(
         "Playlists",
         "Name",
         request.SearchText,
         skip,
         request.PageSize,
         additionalCondition,
         parameters,
         predicateFactory: predicate
      );
      
      var playlistIds = playlists.Items
         .Select(p => p.PlaylistId)
         .ToList();
      
      HashSet<Guid> addedSet = [];

      if (!string.IsNullOrEmpty(videoId) && playlistIds.Count > 0)
      {
         var parseResult = PlatformPrefixToPlatformConverter.ParseVideoPlatform(videoId);

         if (parseResult == null)
         {
            throw new ApiException("Get user playlist error", 400, "Invalid id format");
         }

         (_, actualId) = parseResult.Value;
         addedSet = await _playlistRepository
            .GetPlaylistIdsContainingVideo(actualId, playlistIds);
      }

      var detailedPlaylists = await _playlistRepository.GetPlaylistsDetails(playlistIds);

      var playlistsPreviews = await MapToPreviewDtos(
         detailedPlaylists,
         actualId,
         addedSet,
         requestUserId
      );

      return new PagedResponse<PlaylistPreviewDto>
      {
         Items = playlistsPreviews,
         PageSize = request.PageSize,
         Page = request.Page,
         TotalCount = playlists.Total
      };
   }

   public async Task<PlaylistDto> UpdatePlaylist(Guid requestUserId,UpdatePlaylistRequest request)
   {
      var playlist = await _playlistRepository.FindById(request.PlaylistId);

      if (playlist == null)
      {
         throw new ApiException("Update playlist error", 404, "Playlist wasn't found");
      }

      if (playlist.UserId != requestUserId)
      {
         throw new ApiException("Update playlist error", 403, "You can't update this playlist");
      }
      
      playlist.Name = request.Name;
      playlist.IsPrivate = request.IsPrivate;

      await _playlistRepository.Update(playlist);

      return await MapToPlaylistDto(playlist,requestUserId);
   }

   public async Task DeletePlaylist(Guid userId, Guid playlistId)
   {
      var playlist = await _playlistRepository.FindById(playlistId);

      if (playlist == null)
      {
         throw new ApiException("Delete playlist error", 404, "Playlist wasn't found");
      }

      if (playlist.UserId != userId)
      {
         throw new ApiException("Delete playlist error", 403, "You don't have permission to delete this playlist");
      }

      await _playlistRepository.DeleteAsync(playlistId);
   }

   public async Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId, Guid? requestedUserId)
   {
      var playlist = await _playlistRepository.FindByIdWithVideos(playlistId)
                     ?? throw new ApiException("Get playlist information error", 404, "Playlist wasn't found");

      var playlistDto = await MapToPlaylistDto(playlist, requestedUserId);

      var youtubeIds = playlist.PlaylistVideos
         .Where(pv => pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId))
         .Select(pv => pv.ExternalContentId!)
         .ToList();

      var twitchIds = playlist.PlaylistVideos
         .Where(pv => pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId))
         .Select(pv => pv.ExternalContentId!)
         .ToList();

      var youtubeVideosDict = new Dictionary<string, VideoDto>();
      var twitchStreamsDict = new Dictionary<string, StreamDto>();

      await FetchExternalContentAsync(MediaType.Video, SystemPlatforms.YouTube, youtubeIds, youtubeVideosDict, null);
      await FetchExternalContentAsync(MediaType.LiveStream, SystemPlatforms.Twitch, twitchIds, null, twitchStreamsDict);

      var unavailableExternalCount = (youtubeIds.Count - youtubeVideosDict.Count) + (twitchIds.Count - twitchStreamsDict.Count);

      var visiblePlaylistVideos = playlist.PlaylistVideos
         .Where(pv => 
            (pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId) && youtubeVideosDict.ContainsKey(PlatformPrefixesConstants.YouTubePrefix + pv.ExternalContentId)) ||
            (pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId) && twitchStreamsDict.ContainsKey(PlatformPrefixesConstants.TwitchPrefix + pv.ExternalContentId)) ||
            (pv is { Platform: SystemPlatforms.Webby, Video: not null } && (!pv.Video.IsPrivate || pv.Video.UserId == requestedUserId)))
         .ToList();

      playlistDto.CountOfVideos = visiblePlaylistVideos.Count;

      var sortedVideos = visiblePlaylistVideos.OrderByDescending(pv => pv.CreatedAt).ToList();

      VideoDto? videoDto = null;

      foreach (var playlistVideo in sortedVideos)
      {
         if (playlistVideo is { Platform: SystemPlatforms.Webby, Video: not null })
         {
            var video = playlistVideo.Video;
            videoDto = new VideoDto
            {
               VideoId = PlatformPrefixesConstants.WebbyPrefix + video.VideoId.ToString(),
               Name = video.Name,
               Views = video.Views,
               CreatedAt = video.CreatedAt,
               PreviewUrl = video.PreviewUrl,
               IsPrivate = video.IsPrivate,
            };

            if (video.UserId != Guid.Empty)
            {
               try
               {
                  var userResponse = await _userClient.GetUserByIdAsync(
                     new GetUserRequest
                     {
                        UserId = video.UserId.ToString(),
                        RequestUserId = requestedUserId.ToString()
                     }
                  );

                  if (userResponse != null)
                  {
                     videoDto.User = new UserVideoDto
                     {
                        UserId = userResponse.UserId,
                        Username = userResponse.Username,
                        AvatarUrl = userResponse.AvatarUrl,
                        IsFollowed = userResponse.IsFollowed
                     };
                  }
               }
               catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
               {
                  throw new ApiException("Get playlist information error", 404, ex.Message);
               }
               catch (RpcException ex)
               {
                  throw new ApiException("Get playlist information error", 500, ex.Message);
               }
            }
            break;
         }
         else if (playlistVideo.MediaType == MediaType.Video && playlistVideo.Platform == SystemPlatforms.YouTube && !string.IsNullOrEmpty(playlistVideo.ExternalContentId))
         {
            if (youtubeVideosDict.TryGetValue(PlatformPrefixesConstants.YouTubePrefix + playlistVideo.ExternalContentId, out var ytVideo))
            {
               videoDto = ytVideo;
               break;
            }
         }
         else if (playlistVideo.MediaType == MediaType.LiveStream && playlistVideo.Platform == SystemPlatforms.Twitch && !string.IsNullOrEmpty(playlistVideo.ExternalContentId))
         {
            if (!twitchStreamsDict.TryGetValue(PlatformPrefixesConstants.TwitchPrefix + playlistVideo.ExternalContentId,
                   out var twitchStream)) continue;
            
            videoDto = new VideoDto
            {
               VideoId =  twitchStream.StreamerId,
               Name = twitchStream.Name ?? "Live Stream",
               Views = twitchStream.Viewers,
               CreatedAt = playlistVideo.CreatedAt,
               PreviewUrl = twitchStream.PreviewUrl,
               IsPrivate = false
            };
            break;
         }
      }

      var unavailableWebbyCount = playlist.PlaylistVideos
         .Count(pv => pv is { Platform: SystemPlatforms.Webby, Video.IsPrivate: true } && pv.Video.UserId != requestedUserId);

      return new GetPlaylistResponse
      {
         Playlist = playlistDto,
         FirstVideo = videoDto,
         HiddenVideosCount = unavailableWebbyCount + unavailableExternalCount
      };
   }
   
   public async Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<string> videoIds, Guid requestUserId)
   {
       if (videoIds.Count == 0)
           throw new ApiException("Attach video error", 400, "No videos to add");

       var playlist = await _playlistRepository.GetPlaylistDetails(playlistId)
           ?? throw new ApiException("Attach video to playlist error", 404, "Playlist wasn't found");

       if (playlist.UserId != requestUserId)
           throw new ApiException("Attach video to playlist error", 403, "You can't update this playlist");
       
       var requestedVideos = videoIds
           .Select(id =>
           {
              var result = PlatformPrefixToPlatformConverter.ParseSystemPlatform(id);
              if (result == null) throw new ApiException("Attach video error", 400, "Invalid id format");
              return result.Value;
           })
           .DistinctBy(v => new { v.Platform, v.ActualId })
           .Select(item => new PlaylistVideo
           {
               PlaylistVideoId = Guid.NewGuid(),
               PlaylistId = playlistId,
               Platform = item.Platform,
               InternalContentId = item.Platform == SystemPlatforms.Webby ? Guid.Parse(item.ActualId) : null,
               ExternalContentId = item.Platform != SystemPlatforms.Webby ? item.ActualId : null,
               MediaType = MediaType.Video,
               CreatedAt = DateTime.UtcNow,
           })
           .ToList();
       
       var (itemsToAdd, itemsToDelete) = await _playlistRepository.GetPlaylistItemsDiffAsync(
           playlistId, 
           MediaType.Video, 
           requestedVideos);


       var localVideoGuids = itemsToAdd
           .Where(v => v is { Platform: SystemPlatforms.Webby, InternalContentId: not null })
           .Select(v => v.InternalContentId!.Value)
           .ToList();

       if (localVideoGuids.Any())
       {
           var isAllVideosInDb = await _videoRepository.CheckVideosCount(localVideoGuids);
           if (!isAllVideosInDb)
               throw new ApiException("Update playlist error", 404, "Local videos weren't found");

           var hasForbiddenVideos = await _videoRepository.CheckForbiddenVideos(localVideoGuids, requestUserId);
           if (hasForbiddenVideos)
               throw new ApiException("Add video to playlist error", 403, "You can't add private videos");
       }
       
       if (itemsToAdd.Count > 0)
           await _playlistRepository.AddPlaylistVideos(itemsToAdd);

       if (itemsToDelete.Count > 0)
           await _playlistRepository.DeletePlaylistVideos(itemsToDelete);

       var updatedPlaylist = await _playlistRepository.GetPlaylistDetails(playlistId);

       return await MapToPlaylistDto(updatedPlaylist!, requestUserId);
}
   
   public async Task<PlaylistDto> AttachStreamToPlaylist(Guid playlistId, List<string> streamIds, Guid requestUserId)
   {
      if (streamIds.Count == 0)
         throw new ApiException("Attach stream error", 400, "No streams to add");

      var playlist = await _playlistRepository.GetPlaylistDetails(playlistId)
                     ?? throw new ApiException("Attach stream to playlist error", 404, "Playlist wasn't found");

      if (playlist.UserId != requestUserId)
         throw new ApiException("Attach stream to playlist error", 403, "You can't update this playlist");
      
      var requestedStreams = streamIds
         .Select(id =>
         {
            var result = PlatformPrefixToPlatformConverter.ParseSystemPlatform(id);
            if (result == null) throw new ApiException("Attach stream error", 400, "Invalid id format");
            return result.Value;
         })
         .DistinctBy(s => new { s.Platform, s.ActualId })
         .Select(item => new PlaylistVideo
         {
            PlaylistVideoId = Guid.NewGuid(),
            PlaylistId = playlistId,
            Platform = item.Platform,
            InternalContentId = null, 
            ExternalContentId = item.ActualId,
            MediaType = MediaType.LiveStream,
            CreatedAt = DateTime.UtcNow,
         })
         .ToList();
      
      var (itemsToAdd, itemsToDelete) = await _playlistRepository.GetPlaylistItemsDiffAsync(
         playlistId, 
         MediaType.LiveStream, 
         requestedStreams);

      if (itemsToAdd.Count > 0)
         await _playlistRepository.AddPlaylistVideos(itemsToAdd);

      if (itemsToDelete.Count > 0)
         await _playlistRepository.DeletePlaylistVideos(itemsToDelete);

      var updatedPlaylist = await _playlistRepository.GetPlaylistDetails(playlistId);

      return await MapToPlaylistDto(updatedPlaylist!, requestUserId);
   }

   public async Task<PagedResponse<SearchPlaylistDto>> SearchPlaylists(
      Guid? requestUserId,
      SearchOptions searchOptions)
   {
      var skip = (searchOptions.Page - 1) * searchOptions.PageSize;

      var (additionalCondition, parameters, predicate) =
         PlaylistSearchFilter.SearchPlaylistFilters(requestUserId);

      var (playlists, totalPlaylists) = await _playlistRepository.SearchAsync(
         "Playlists",
         "Name",
         searchOptions.SearchText,
         skip,
         searchOptions.PageSize,
         additionalCondition,
         parameters,
         predicateFactory: predicate
      );

      var playlistIds = playlists.Select(p => p.PlaylistId).ToList();
      var detailedPlaylists = await _playlistRepository.GetPlaylistsDetails(playlistIds);

      var userIds = detailedPlaylists
         .Select(p => p.UserId.ToString())
         .Distinct()
         .ToList();

      var users = await _userClient.GetUsersByIdsAsync(new GetUsersRequest
         { UserIds = {userIds}}
      );


      var usersDict = users.Users.ToDictionary(u => u.UserId, u => u);
         
      var tasks = detailedPlaylists
         .Select(p => MapToSearchPlaylistDto(p, requestUserId, usersDict));
      
      var items = (await Task.WhenAll(tasks)).ToList();

      return new PagedResponse<SearchPlaylistDto>
      {
         Items = items,
         TotalCount = totalPlaylists,
         Page = searchOptions.Page,
         PageSize = searchOptions.PageSize,
      };
   }

   public async Task<bool> CheckIfVideoExistInPlaylist(Guid playlistId, string videoId)
   {
      var parseResult = PlatformPrefixToPlatformConverter.ParseVideoPlatform(videoId);

      if (parseResult == null)
      {
         throw new ApiException("Check if video exist error", 400, "Invalid id format type");
      }

      var (_, actualId) = parseResult.Value;
      return await _playlistRepository.CheckIsVideoAdded(actualId, playlistId);
   }
      

   private async Task<PlaylistDto> MapToPlaylistDto(Playlist playlist, Guid? requestUserId)
   {
      var covers = await GetPlaylistsCoversAsync(new[] { playlist }, requestUserId);
      var coverUrl = covers.GetValueOrDefault(playlist.PlaylistId, DefaultLinks.PlaylistEmptyLink);

      return new PlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = coverUrl
      };
   }
   
   private async Task<SearchPlaylistDto> MapToSearchPlaylistDto(
      Playlist playlist,
      Guid? requestUserId,
      Dictionary<string, UserResponse> usersDict)
   {
      var covers = await GetPlaylistsCoversAsync(new[] { playlist }, requestUserId);
      var coverUrl = covers.GetValueOrDefault(playlist.PlaylistId, DefaultLinks.PlaylistEmptyLink);

      usersDict.TryGetValue(playlist.UserId.ToString(), out var user);

      return new SearchPlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         Username = user?.Username ?? "Deleted user",
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = coverUrl
      };
   }
   
   private async Task<List<PlaylistPreviewDto>> MapToPreviewDtos(
      IEnumerable<Playlist> playlists,
      string videoId,
      HashSet<Guid> addedSet,
      Guid? requestUserId)
   {
      var covers = await GetPlaylistsCoversAsync(playlists, requestUserId);
      
      return playlists
         .Select(p => new PlaylistPreviewDto
         {
            PlaylistId = p.PlaylistId,
            Name = p.Name,
            CountOfVideos = p.PlaylistVideos?.Count ?? 0,
            PlaylistCover = covers.GetValueOrDefault(p.PlaylistId, DefaultLinks.PlaylistEmptyLink),
            IsVideoAdded = !string.IsNullOrEmpty(videoId) && addedSet.Contains(p.PlaylistId),
            IsPrivate = p.IsPrivate
         })
         .ToList();
   }
   
   
   private async Task<Dictionary<Guid, string>> GetPlaylistsCoversAsync(IEnumerable<Playlist> playlists, Guid? requestUserId)
   {
       var covers = new Dictionary<Guid, string>();
       var youtubeIdsToFetch = new HashSet<string>();
       var twitchIdsToFetch = new HashSet<string>();
       
       foreach (var playlist in playlists)
       {
           if (playlist.PlaylistVideos == null) continue;

           var youtubeVideos = playlist.PlaylistVideos
               .Where(pv => pv.MediaType == MediaType.Video && pv.Platform == SystemPlatforms.YouTube && !string.IsNullOrEmpty(pv.ExternalContentId))
               .Select(pv => pv.ExternalContentId!);

           foreach (var id in youtubeVideos)
           {
               youtubeIdsToFetch.Add(id);
           }

           var twitchStreams = playlist.PlaylistVideos
               .Where(pv => pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId))
               .Select(pv => pv.ExternalContentId!);

           foreach (var id in twitchStreams)
           {
               twitchIdsToFetch.Add(id);
           }
       }
       
       var youtubeVideosDict = new Dictionary<string, VideoDto>();
       if (youtubeIdsToFetch.Count != 0)
       {
           var ytList = await _youtubeSearchService.GetList(youtubeIdsToFetch.ToList());
           youtubeVideosDict = ytList.ToDictionary(v => v.VideoId!);
       }

       var twitchStreamsDict = new Dictionary<string, StreamDto>();
       if (twitchIdsToFetch.Count != 0)
       {
           var twitchList = await _twitchSearchService.GetList(twitchIdsToFetch.ToList());
           twitchStreamsDict = twitchList.ToDictionary(s => s.StreamerId!);
       }
       
       foreach (var playlist in playlists)
       {
           string coverUrl = DefaultLinks.PlaylistEmptyLink;

           var availableItems = playlist.PlaylistVideos?
               .Where(pv =>
                   (pv.MediaType == MediaType.Video && pv.Platform == SystemPlatforms.YouTube) ||
                   (pv.MediaType == MediaType.LiveStream && pv.Platform == SystemPlatforms.Twitch) ||
                   (pv.MediaType == MediaType.Video && pv.Platform == SystemPlatforms.Webby && pv.Video != null && (!pv.Video.IsPrivate || pv.Video.UserId == requestUserId))
               )
               .OrderByDescending(pv => pv.CreatedAt)
               .ToList();

           if (availableItems != null && availableItems.Any())
           {
               foreach (var pv in availableItems)
               {
                   if (pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.Webby })
                   {
                       var localCover = pv.Video!.PreviewUrl;
                       if (!string.IsNullOrEmpty(localCover))
                       {
                           coverUrl = localCover;
                           break;
                       }
                   }
                   else if (pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId))
                   {
                       if (youtubeVideosDict.TryGetValue(PlatformPrefixesConstants.YouTubePrefix + pv.ExternalContentId, out var ytVideo) &&
                           !string.IsNullOrEmpty(ytVideo.PreviewUrl))
                       {
                           coverUrl = ytVideo.PreviewUrl;
                           break;
                       }
                   }
                   else if (pv.MediaType == MediaType.LiveStream && pv.Platform == SystemPlatforms.Twitch && !string.IsNullOrEmpty(pv.ExternalContentId))
                   {
                       if (twitchStreamsDict.TryGetValue(PlatformPrefixesConstants.TwitchPrefix + pv.ExternalContentId, out var twitchStream) &&
                           !string.IsNullOrEmpty(twitchStream.PreviewUrl))
                       {
                           coverUrl = twitchStream.PreviewUrl;
                           break;
                       }
                   }
               }
           }

           covers[playlist.PlaylistId] = coverUrl;
       }

       return covers;
   }
   
   private async Task FetchExternalContentAsync(
      MediaType mediaType, 
      SystemPlatforms platform, 
      List<string> ids, 
      Dictionary<string, VideoDto>? youtubeDict, 
      Dictionary<string, StreamDto>? twitchDict)
   {
      if (ids.Count == 0) return;

      switch (mediaType, platform)
      {
         case (MediaType.Video, SystemPlatforms.YouTube):
            var ytList = await _youtubeSearchService.GetList(ids);
            if (youtubeDict != null)
            {
               foreach (var item in ytList)
                  youtubeDict[item.VideoId!] = item;
            }
            break;

         case (MediaType.LiveStream, SystemPlatforms.Twitch):
            var twitchList = await _twitchSearchService.GetList(ids);
            if (twitchDict != null)
            {
               foreach (var item in twitchList)
                  twitchDict[item.StreamerId!] = item;
            }
            break;
      }
   }
   
}