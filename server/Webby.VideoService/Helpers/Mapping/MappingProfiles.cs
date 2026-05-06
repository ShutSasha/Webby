using AutoMapper;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Models;

namespace Webby.VideoService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<Models.Playlist, PlaylistDto>();
      CreateMap<Models.Video, VideoDto>()
         .ForMember(dest 
            => dest.VideoId,
            opt => opt.MapFrom(src => $"{PlatformPrefixesConstants.WebbyPrefix}{src.VideoId}"))
         .ForMember(dest => dest.VideoTags, opt => opt.MapFrom(src => 
            src.VideoTags
               .Select(vt => vt.Tag.Name)
               .ToList()));
      
      CreateMap<Models.Video, UploadVideoResponse>();
      CreateMap<Models.Video, PreviewVideoDto>()
         .ForMember(dest
               => dest.VideoId,
            opt =>
               opt.MapFrom(src => $"{PlatformPrefixesConstants.WebbyPrefix}{src.VideoId}"));
   }
}