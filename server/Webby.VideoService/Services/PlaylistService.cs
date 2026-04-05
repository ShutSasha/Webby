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

namespace Webby.VideoService.Services;

public class PlaylistService : IPlaylistService
{
   private readonly IPlaylistRepository _playlistRepository;
   private readonly IMapper _mapper;
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   private readonly IVideoRepository _videoRepository;

   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper,
      UserGrpcService.UserGrpcServiceClient userClient, IVideoRepository videoRepository)
   {
      _playlistRepository = playlistRepository;
      _mapper = mapper;
      _userClient = userClient;
      _videoRepository = videoRepository;
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

   public async Task<PagedResponse<PlaylistPreviewDto>> GetUserPlaylists(Guid? requestUserId, Guid? videoId, Guid userId, GetUserPlaylistsRequest request)
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

      if (videoId.HasValue && playlistIds.Count > 0)
      {
         addedSet = await _playlistRepository
            .GetPlaylistIdsContainingVideo(videoId.Value, playlistIds);
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

      return MapToPlaylistDto(playlist,requestUserId);
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

      var playlistDto = MapToPlaylistDto(playlist,requestedUserId);

      var visiblePlaylistVideos = playlist.PlaylistVideos
         .Where(pv => !pv.Video.IsPrivate || pv.Video.UserId == requestedUserId)
         .ToList();

      playlistDto.CountOfVideos = visiblePlaylistVideos.Count;

      var video = visiblePlaylistVideos
         .OrderByDescending(pv => pv.CreatedAt)
         .Select(pv => pv.Video)
         .FirstOrDefault();

      VideoDto? videoDto = null;

      if (video != null)
      {
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

      var unavailableCount = playlist.PlaylistVideos
         .Select(pv => pv.Video)
         .Count(v => v.IsPrivate && v.UserId != requestedUserId);
      
      return new GetPlaylistResponse
      {
         Playlist = playlistDto,
         FirstVideo = videoDto,
         HiddenVideosCount = unavailableCount
      };
   }

   public async Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<Guid> videoIds, Guid requestUserId)
   {
      if (videoIds == null || !videoIds.Any())
         throw new ApiException("Attach video error", 400, "No videos to add");

      var playlist = await _playlistRepository.GetPlaylistDetails(playlistId)
                     ?? throw new ApiException("Attach video to playlist error", 404, "Playlist wasn't found");

      if (playlist.UserId != requestUserId)
      {
         throw new ApiException("Attach video to playlist error", 403, "You can't update this playlist");
      }

      var existingVideoIds = playlist.PlaylistVideos
         .Select(pv => pv.VideoId)
         .ToHashSet();
      
      var playlistVideosToDelete = playlist.PlaylistVideos
         .Where(p => videoIds.Contains(p.VideoId))
         .ToList();

      var playlistVideos = videoIds
         .Distinct()
         .Where(videoId => !existingVideoIds.Contains(videoId))
         .Select(videoId => new PlaylistVideo
         {
            PlaylistId = playlistId,
            VideoId = videoId,
            CreatedAt = DateTime.UtcNow,
         })
         .ToList();

      var playlistVideoIds = playlistVideos.Select(pv => pv.VideoId).ToList();

      if (playlistVideoIds.Count > 0)
      {
         var isAllVideosInDb = await _videoRepository.CheckVideosCount(playlistVideoIds);
         if (!isAllVideosInDb)
         {
            throw new ApiException("Update playlist error", 404, "Videos wasn't found");
         }

         var hasForbiddenVideos = await _videoRepository.CheckForbiddenVideos(playlistVideoIds, requestUserId);
         
         if (hasForbiddenVideos)
         {
            throw new ApiException("Add video to playlist error", 403, "You can't add private videos");
         }
      }

      await _playlistRepository.AddPlaylistVideos(playlistVideos);
      await _playlistRepository.DeletePlaylistVideos(playlistVideosToDelete);
      return MapToPlaylistDto((await _playlistRepository.GetPlaylistDetails(playlistId))!,requestUserId);
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
         { UserIds = {userIds }
   }
   );


   var usersDict = users.Users.ToDictionary(u => u.UserId, u => u);
      
      var items = detailedPlaylists
         .Select(p => MapToSearchPlaylistDto(p,requestUserId, usersDict))
         .ToList();

      return new PagedResponse<SearchPlaylistDto>
      {
         Items = items,
         TotalCount = totalPlaylists,
         Page = searchOptions.Page,
         PageSize = searchOptions.PageSize,
      };
   }

   public async Task<bool> CheckIfVideoExistInPlaylist(Guid playlistId, Guid videoId)
      => await _playlistRepository.CheckIsVideoAdded(videoId,playlistId);

   //TODO: fix reusing search
   public async Task<GlobalSearchPlaylistResponse> GlobalPlaylistsSearch(Guid? requestUserId, GlobalSearchOptions searchOptions)
   {
      if (searchOptions.SectionType != SearchSections.Playlists)
      {
         throw new ApiException("Global playlists search", 400, "Incorrect search section type");
      }
      
      var playlistsResult = await SearchPlaylists(requestUserId, new SearchOptions()
      {
         Page = searchOptions.Page,
         PageSize = 5,
         SearchText = searchOptions.SearchText
      });
      
      return new GlobalSearchPlaylistResponse()
      {
         WebbyPlaylists = new SearchSection<SearchPlaylistDto>()
         {
            Items = playlistsResult.Items,
            Page =  playlistsResult.Page,
            TotalCount = playlistsResult.TotalCount
         }
      };
   }

   private PlaylistDto MapToPlaylistDto(Playlist playlist,Guid? requestUserId)
   {
      var lastVideo = playlist.PlaylistVideos?
         .OrderByDescending(pv => pv.CreatedAt)
         .Select(pv => pv.Video!)
         .FirstOrDefault(v => !v.IsPrivate || v.UserId == requestUserId);

      return new PlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
   
   private SearchPlaylistDto MapToSearchPlaylistDto(
      Playlist playlist,
      Guid? requestUserId,
      Dictionary<string, UserResponse> usersDict)
   {
      var lastVideo = playlist.PlaylistVideos?
         .OrderByDescending(pv => pv.CreatedAt)
         .Select(pv => pv.Video!)
         .FirstOrDefault(v => !v.IsPrivate || v.UserId == requestUserId);

      usersDict.TryGetValue(playlist.UserId.ToString(), out var user);

      return new SearchPlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         Username = user?.Username ?? "Deleted user",
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
   
   private List<PlaylistPreviewDto> MapToPreviewDtos(
      IEnumerable<Playlist> playlists,
      Guid? videoId,
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
               IsVideoAdded = videoId.HasValue && addedSet.Contains(p.PlaylistId),
               IsPrivate = p.IsPrivate
            };
         })
         .ToList();
   }
}