using AutoMapper;
using Grpc.Core;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;

namespace Webby.VideoService.Services;

public class PlaylistService : IPlaylistService
{
   private readonly IPlaylistRepository _playlistRepository;
   private readonly IMapper _mapper;
   private readonly UserGrpcService.UserGrpcServiceClient _userClient;
   
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
         Name = request.Name,
         UserId = userId
      };

      await _playlistRepository.Add(playlist);

      return _mapper.Map<PlaylistDto>(playlist);
   }

   public async Task<List<PlaylistDto>> GetUserPlaylists(Guid userId)
   {
      var playlists = await _playlistRepository.GetUserPlaylists(userId);

      return playlists
         .Select(MapToPlaylistDto)
         .ToList();
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
   
   public async Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId)
{
    var playlist = await _playlistRepository.FindByIdWithVideos(playlistId);

    if (playlist == null)
    {
        throw new ApiException("Get playlist information error", 404, "Playlist wasn't found");
    }
    
    var playlistDto = MapToPlaylistDto(playlist);
    
    var videos = playlist.PlaylistVideos
        .Select(pv => pv.Video)
        .OrderByDescending(v => v.CreatedAt)
        .ToList();

    var videoDtos = videos.Select(v => new VideoDto
    {
        VideoId = v.VideoId,
        Name = v.Name,
        Views = v.Views,
        CreatedAt = v.CreatedAt,
        PreviewUrl = v.PreviewUrl,
        IsPrivate = v.IsPrivate,
        User = null 
    }).ToList();

    var userIds = videos
       .Select(v => v.UserId.ToString())
       .Distinct()
       .ToList();
    

    if (userIds.Any())
    {
         var usersResponse = await _userClient.GetUsersByIdsAsync(
             new GetUsersRequest { UserIds = { userIds } }
         );

         var usersDict = usersResponse.Users
             .ToDictionary(
                 u => Guid.Parse(u.UserId),
                 u => new UserVideoDto
                 {
                     UserId = Guid.Parse(u.UserId),
                     Username = u.Username,
                     AvatarUrl = u.AvatarUrl
                 });
         
         foreach (var videoDto in videoDtos)
         {
             var originalVideo = videos.First(v => v.VideoId == videoDto.VideoId);

             if (usersDict.TryGetValue(originalVideo.UserId, out var user))
             {
                 videoDto.User = user;
             }
         }
    }
    
    return new GetPlaylistResponse
    {
        Playlist = playlistDto,
        Videos = videoDtos
    };
}
   
   public async Task<PlaylistDto> AttachVideoToPlaylist(Guid playlistId, List<Guid> videoIds)
   {
      if (videoIds == null || !videoIds.Any())
         throw new ApiException("Attach video error", 400, "No videos to add");

      var playlist = await _playlistRepository.FindById(playlistId);

      if (playlist == null)
         throw new ApiException("Attach to playlist error", 404,"Playlist wasn't found");

      var playlistVideos = videoIds
         .Distinct()
         .Select(videoId => new PlaylistVideo
         {
            PlaylistId = playlistId,
            VideoId = videoId
         })
         .ToList();

      await _playlistRepository.AddPlaylistVideos(playlistVideos);
      return MapToPlaylistDto(playlist);
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
         Description = playlist.Description,
         CountOfVideos = playlist.PlaylistVideos?.Count ?? 0,
         PlaylistCover = lastVideo?.Video.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
      };
   }
}