using AutoMapper;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Models;

namespace Webby.VideoService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<Models.Playlist, PlaylistDto>();
      CreateMap<Models.Video, VideoDto>();
   }
}