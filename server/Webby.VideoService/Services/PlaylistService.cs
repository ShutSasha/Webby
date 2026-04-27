using System.Net.NetworkInformation;
using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
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
   private readonly YoutubeSearchService _youtubeSearchService;
   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper,
      UserGrpcService.UserGrpcServiceClient userClient, IVideoRepository videoRepository, YoutubeSearchService youtubeSearchService)
   {
      _playlistRepository = playlistRepository;
      _mapper = mapper;
      _userClient = userClient;
      _videoRepository = videoRepository;
      _youtubeSearchService = youtubeSearchService;
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

      if (videoId != string.Empty && playlistIds.Count > 0)
      {
         addedSet = await _playlistRepository
            .GetPlaylistIdsContainingVideo(videoId, playlistIds);
      }

      var detailedPlaylists = await _playlistRepository.GetPlaylistsDetails(playlistIds);

      var playlistsPreviews = MapToPreviewDtos(
         detailedPlaylists,
         videoId,
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


      var visiblePlaylistVideos = playlist.PlaylistVideos
         .Where(pv => pv.VideoPlatform == VideoPlatform.YouTube || 
                     (pv is { VideoPlatform: VideoPlatform.Webby, Video: not null } && (!pv.Video.IsPrivate || pv.Video.UserId == requestedUserId)))
         .ToList();

      playlistDto.CountOfVideos = visiblePlaylistVideos.Count;
      
      var latestPlaylistVideo = visiblePlaylistVideos.MaxBy(pv => pv.CreatedAt);

      VideoDto? videoDto = null;

      if (latestPlaylistVideo != null)
      {
         if (latestPlaylistVideo is { VideoPlatform: VideoPlatform.Webby, Video: not null })
         {
            var video = latestPlaylistVideo.Video;
            videoDto = new VideoDto
            {
               VideoId = video.VideoId.ToString(),
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
         }
         else if (latestPlaylistVideo.VideoPlatform == VideoPlatform.YouTube && !string.IsNullOrEmpty(latestPlaylistVideo.ExternalVideoId))
         {
            var youtubeVideos = await _youtubeSearchService.FindById(latestPlaylistVideo.ExternalVideoId);
            videoDto = youtubeVideos;
         }
      }
      
      var unavailableCount = playlist.PlaylistVideos
         .Count(pv => pv is { VideoPlatform: VideoPlatform.Webby, Video.IsPrivate: true } && pv.Video.UserId != requestedUserId);
      
      return new GetPlaylistResponse
      {
         Playlist = playlistDto,
         FirstVideo = videoDto,
         HiddenVideosCount = unavailableCount
      };
   }
   
   public async Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<AddVideoToPlaylistItem> videoItems, Guid requestUserId)
   {
       if (videoItems == null || !videoItems.Any())
           throw new ApiException("Attach video error", 400, "No videos to add");

       var playlist = await _playlistRepository.GetPlaylistDetails(playlistId)
                      ?? throw new ApiException("Attach video to playlist error", 404, "Playlist wasn't found");

       if (playlist.UserId != requestUserId)
       {
           throw new ApiException("Attach video to playlist error", 403, "You can't update this playlist");
       }

       var uniqueRequestedItems = videoItems
           .DistinctBy(v => new { v.VideoPlatform, v.ItemId })
           .ToList();
       
       var existingItems = playlist.PlaylistVideos
           .Select(pv => new 
           { 
               pv.VideoPlatform, 
               ItemId = pv.VideoPlatform == VideoPlatform.Webby ? pv.VideoId.ToString() : pv.ExternalVideoId 
           })
           .ToHashSet();
       
       var playlistVideosToDelete = playlist.PlaylistVideos
           .Where(p => uniqueRequestedItems.Any(req => 
               req.VideoPlatform == p.VideoPlatform && 
               req.ItemId == (p.VideoPlatform == VideoPlatform.Webby ? p.VideoId.ToString() : p.ExternalVideoId)))
           .ToList();
       
       var newItemsToAdd = uniqueRequestedItems
          .Where(req => !existingItems.Contains(new { req.VideoPlatform, req.ItemId }))
          .ToList();

       var localVideoIdsStrings = newItemsToAdd
           .Where(v => v.VideoPlatform == VideoPlatform.Webby)
           .Select(v => v.ItemId)
           .ToList();

       if (localVideoIdsStrings.Any())
       {
           var localVideoGuids = localVideoIdsStrings
               .Select(id => Guid.TryParse(id, out var guid) ? guid : Guid.Empty)
               .Where(g => g != Guid.Empty)
               .ToList();

           if (localVideoGuids.Count != localVideoIdsStrings.Count)
           {
               throw new ApiException("Update playlist error", 400, "Invalid local video ID format");
           }

           var isAllVideosInDb = await _videoRepository.CheckVideosCount(localVideoGuids);
           if (!isAllVideosInDb)
           {
               throw new ApiException("Update playlist error", 404, "Local videos weren't found");
           }

           var hasForbiddenVideos = await _videoRepository.CheckForbiddenVideos(localVideoGuids, requestUserId);
           if (hasForbiddenVideos)
           {
               throw new ApiException("Add video to playlist error", 403, "You can't add private videos");
           }
       }
       
       var playlistVideos = newItemsToAdd.Select(item => new PlaylistVideo
       {
          PlaylistVideoId = Guid.NewGuid(),
           PlaylistId = playlistId,
           VideoPlatform = item.VideoPlatform,
           VideoId= item.VideoPlatform == VideoPlatform.Webby ? Guid.Parse(item.ItemId) : null,
           ExternalVideoId = item.VideoPlatform == VideoPlatform.YouTube ? item.ItemId : null,
           CreatedAt = DateTime.UtcNow,
       }).ToList();
       
       if (playlistVideos.Count > 0)
       {
           await _playlistRepository.AddPlaylistVideos(playlistVideos);
       }

       if (playlistVideosToDelete.Count > 0)
       {
           await _playlistRepository.DeletePlaylistVideos(playlistVideosToDelete);
       }

       return await MapToPlaylistDto((await _playlistRepository.GetPlaylistDetails(playlistId))!, requestUserId);
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
      => await _playlistRepository.CheckIsVideoAdded(videoId, playlistId);

   private async Task<PlaylistDto> MapToPlaylistDto(Playlist playlist, Guid? requestUserId)
   {
      string coverUrl = DefaultLinks.PlaylistEmptyLink;
      
      var lastVideoEntry = playlist.PlaylistVideos?
         .OrderByDescending(pv => pv.CreatedAt)
         .FirstOrDefault(pv => 
            pv.VideoPlatform == VideoPlatform.YouTube || 
            (pv is { VideoPlatform: VideoPlatform.Webby, Video: not null } && (!pv.Video.IsPrivate || pv.Video.UserId == requestUserId))
         );

      if (lastVideoEntry != null)
      {
         if (lastVideoEntry.VideoPlatform == VideoPlatform.Webby)
         {
            coverUrl = lastVideoEntry.Video!.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink;
         }
         else if (lastVideoEntry.VideoPlatform == VideoPlatform.YouTube)
         {
            try 
            {
               var youtubeVideo = await _youtubeSearchService.FindById(lastVideoEntry.ExternalVideoId!);
               coverUrl = youtubeVideo?.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink;
            }
            catch
            {
               coverUrl = DefaultLinks.PlaylistEmptyLink;
            }
         }
      }

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
      string lastVideoCoverThumbnail = DefaultLinks.PlaylistEmptyLink;
      
      var lastVideoEntry = playlist.PlaylistVideos?
         .OrderByDescending(pv => pv.CreatedAt)
         .FirstOrDefault(pv =>
            pv.VideoPlatform == VideoPlatform.YouTube ||
            pv is {VideoPlatform: VideoPlatform.Webby, Video: not null, Video.IsPrivate: false }
            );

      if (lastVideoEntry != null)
      {
         if (lastVideoEntry.VideoPlatform == VideoPlatform.Webby)
         {
            lastVideoCoverThumbnail = lastVideoEntry.Video!.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink;
         }
         else if (lastVideoEntry.VideoPlatform == VideoPlatform.YouTube)
         {
            try
            {
               var youtubeVideo = await _youtubeSearchService.FindById(lastVideoEntry.ExternalVideoId!);
               lastVideoCoverThumbnail = youtubeVideo.PreviewUrl;
            }
            catch
            {
               lastVideoCoverThumbnail = DefaultLinks.PlaylistEmptyLink;
            }

         }
      }

      usersDict.TryGetValue(playlist.UserId.ToString(), out var user);

      return new SearchPlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         Username = user?.Username ?? "Deleted user",
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideoCoverThumbnail
      };
   }
   
   private List<PlaylistPreviewDto> MapToPreviewDtos(
      IEnumerable<Playlist> playlists,
      string videoId,
      HashSet<Guid> addedSet,
      Guid? requestUserId)
   {
      return playlists
         .Select(p => 
         {
            var lastVideo = p.PlaylistVideos
               .OrderByDescending(pv => pv.CreatedAt)
               .Select(pv => pv.Video!)
               .FirstOrDefault(v => !v.IsPrivate || v.UserId == requestUserId);

            return new PlaylistPreviewDto
            {
               PlaylistId = p.PlaylistId,
               Name = p.Name,
               CountOfVideos = p.PlaylistVideos.Count,
               PlaylistCover = lastVideo?.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink,
               IsVideoAdded = videoId != string.Empty && addedSet.Contains(p.PlaylistId),
               IsPrivate = p.IsPrivate
            };
         })
         .ToList();
   }
}