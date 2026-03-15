using AutoMapper;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;
using Webby.VideoService.Repositories;

namespace Webby.VideoService.Services;

public class PlaylistService : IPlaylistService
{
   private readonly IPlaylistRepository _playlistRepository;
   private readonly IMapper _mapper;

   public PlaylistService(IPlaylistRepository playlistRepository, IMapper mapper)
   {
      _playlistRepository = playlistRepository;
      _mapper = mapper;
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

      return playlists.Select(p =>
      {
         var lastVideo = p.PlaylistVideos
            ?.OrderByDescending(pv => pv.Video.CreatedAt)
            .FirstOrDefault();

         return new PlaylistDto
         {
            PlaylistId = p.PlaylistId,
            UserId = p.UserId,
            Name = p.Name,
            Description = p.Description,
            CountOfVideos = p.PlaylistVideos?.Count ?? 0,
            PlaylistCover = lastVideo?.Video.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink
         };
      }).ToList();
      
   }

   public Task<PlaylistDto> UpdatePlaylist(UpdatePlaylistRequest request)
   {
      throw new NotImplementedException();
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

   public Task<GetPlaylistResponse> GetPlaylistInformation(Guid playlistId)
   {
      throw new NotImplementedException();
   }
}