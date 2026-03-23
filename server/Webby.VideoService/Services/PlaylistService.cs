using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
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
   
   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper, UserGrpcService.UserGrpcServiceClient userClient)
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
         Description = request.Description,
         IsPrivate = request.IsPrivate,
         CreatedAt = DateTime.UtcNow,
         Name = request.Name,
         UserId = userId
      };

      await _playlistRepository.Add(playlist);

      return _mapper.Map<PlaylistDto>(playlist);
   }

   public async Task<PagedResponse<PlaylistDto>> GetUserPlaylists(Guid? requestUserId, Guid userId, int page, int pageSize)
   {
      page = page <= 0 ? 1 : page;
      pageSize = pageSize <= 0 ? 10 : pageSize;

      var (playlists, totalCount) = await _playlistRepository
         .GetPaginatedUserPlaylists(requestUserId == userId,userId, page, pageSize);

      return new PagedResponse<PlaylistDto>
      {
         Items = playlists.Select(MapToPlaylistDto).ToList(),
         Page = page,
         PageSize = pageSize,
         TotalCount = totalCount
      };

   }

   public async Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request)
   {
      var playlist = await _playlistRepository.FindById(request.PlaylistId);

      if (playlist == null)
      {
         throw new ApiException("Update playlist error", 404, "Playlist wasn't found");
      }

      playlist.Description = request.Description;
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
         throw new ApiException("Attach to playlist error", 404,"Playlist wasn't found");

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

   public async Task<PagedResponse<PlaylistDto>> SearchPlaylists(SearchOptions searchOptions)
   {
      var skip = (searchOptions.Page - 1) * searchOptions.PageSize;

      var (playlists, totalPlaylists) = await _playlistRepository
         .SearchPlaylistsAsync(searchOptions.SearchText, skip, searchOptions.PageSize);

      var items = playlists.Select(MapToPlaylistDto).ToList();

      return new PagedResponse<PlaylistDto>
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
         Description = playlist.Description ?? string.Empty,
         IsPrivate = playlist.IsPrivate,
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.Video.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
}