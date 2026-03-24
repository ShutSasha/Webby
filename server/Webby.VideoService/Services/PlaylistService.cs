using System.Net.NetworkInformation;
using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
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
   private readonly IVideoService _videoService;

   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper,
      UserGrpcService.UserGrpcServiceClient userClient)
   {
      _playlistRepository = playlistRepository;
      _mapper = mapper;
      _userClient = userClient;
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

   public async Task<PagedResponse<PlaylistPreviewDto>> GetUserPlaylists
   (
      Guid? requestUserId,
      Guid? videoId,
      Guid userId,
      SearchOptions searchOptions
   )
   {
      var skip = (searchOptions.Page - 1) * searchOptions.PageSize;

      var (additionalCondition, parameters, predicate)
         = PlaylistSearchFilter.SearchUserPlaylistsFilter(userId, requestUserId == userId);

      var playlists = await _playlistRepository.SearchAsync(
         "Playlists",
         "Name",
         searchOptions.SearchText,
         skip,
         searchOptions.PageSize,
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
         addedSet
      );

      return new PagedResponse<PlaylistPreviewDto>
      {
         Items = playlistsPreviews,
         PageSize = searchOptions.PageSize,
         Page = searchOptions.Page,
         TotalCount = playlists.Total
      };
   }

   public async Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request)
   {
      var playlist = await _playlistRepository.FindById(request.PlaylistId);

      if (playlist == null)
      {
         throw new ApiException("Update playlist error", 404, "Playlist wasn't found");
      }
      
      playlist.Name = request.Name;
      playlist.IsPrivate = request.IsPrivate;

      await _playlistRepository.Update(playlist);

      return MapToPlaylistDto(playlist);
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

      var playlistDto = MapToPlaylistDto(playlist);

      var video = playlist.PlaylistVideos
         .Select(pv => pv.Video)
         .MaxBy(v => v.CreatedAt);

      VideoDto? videoDto = null;

      if (video != null)
      {
         videoDto = new VideoDto
         {
            VideoId = video.VideoId,
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
                     UserId = Guid.Parse(userResponse.UserId),
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

      return new GetPlaylistResponse
      {
         Playlist = playlistDto,
         FirstVideo = videoDto
      };
   }

   public async Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<Guid> videoIds)
   {
      if (videoIds == null || !videoIds.Any())
         throw new ApiException("Attach video error", 400, "No videos to add");

      var playlist = await _playlistRepository.GetPlaylistDetails(playlistId);

      if (playlist == null)
         throw new ApiException("Attach to playlist error", 404, "Playlist wasn't found");

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
            VideoId = videoId
         })
         .ToList();

      await _playlistRepository.AddPlaylistVideos(playlistVideos);
      await _playlistRepository.DeletePlaylistVideos(playlistVideosToDelete);
      return MapToPlaylistDto((await _playlistRepository.GetPlaylistDetails(playlistId))!);
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
         .Select(p => MapToSearchPlaylistDto(p, usersDict))
         .ToList();

      return new PagedResponse<SearchPlaylistDto>
      {
         Items = items,
         TotalCount = totalPlaylists,
         Page = searchOptions.Page,
         PageSize = searchOptions.PageSize,
      };
   }

   private PlaylistDto MapToPlaylistDto(Playlist playlist)
   {
      var lastVideo = playlist.PlaylistVideos?
         .MaxBy(pv => pv.Video.CreatedAt);

      return new PlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.Video.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
   
   private SearchPlaylistDto MapToSearchPlaylistDto(
      Playlist playlist,
      Dictionary<string, UserResponse> usersDict)
   {
      var lastVideo = playlist.PlaylistVideos?
         .MaxBy(pv => pv.Video.CreatedAt);
      
      usersDict.TryGetValue(playlist.UserId.ToString(), out var user);

      return new SearchPlaylistDto
      {
         PlaylistId = playlist.PlaylistId,
         UserId = playlist.UserId,
         Name = playlist.Name,
         IsPrivate = playlist.IsPrivate,
         Username = user?.Username ?? "Deleted user",
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.Video.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
   
   private List<PlaylistPreviewDto> MapToPreviewDtos(
      IEnumerable<Playlist> playlists,
      Guid? videoId,
      HashSet<Guid> addedSet)
   {
      return playlists
         .Select(p => new PlaylistPreviewDto
         {
            PlaylistId = p.PlaylistId,
            Name = p.Name,
            CountOfVideos = p.PlaylistVideos.Count,
            PlaylistCover = p.PlaylistVideos
                               .MaxBy(pv => pv.Video.CreatedAt)?.Video.PreviewUrl 
                            ?? DefaultLinks.PlaylistEmptyLink,
            IsVideoAdded = videoId.HasValue && addedSet.Contains(p.PlaylistId),
            IsPrivate = p.IsPrivate
         })
         .ToList();
   }
}