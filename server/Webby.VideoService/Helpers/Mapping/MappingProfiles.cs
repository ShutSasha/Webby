using AutoMapper;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<Models.Playlist, PlaylistDto>();

      CreateMap<Models.Video, VideoDto>()
         .ForMember(dest => dest.VideoId,
            opt => opt.MapFrom(src => $"{PlatformPrefixesConstants.WebbyPrefix}{src.VideoId}"))
         .ForMember(dest => dest.VideoTags, opt => opt.MapFrom(src => src.VideoTags.Select(vt => vt.Tag.Name).ToList()))
         .ForMember(dest => dest.MediaType, opt => opt.MapFrom(src => MediaType.Video));
      
         
      CreateMap<Models.Video, UploadVideoResponse>()
         .ForMember(dest => dest.VideoId, opt => opt.MapFrom(src => $"{PlatformPrefixesConstants.WebbyPrefix}{src.VideoId}"));
      
      CreateMap<Models.Video, PreviewVideoDto>()
         .ForMember(dest => dest.VideoId, opt => opt.MapFrom(src => $"{PlatformPrefixesConstants.WebbyPrefix}{src.VideoId}"));

      CreateMap<StreamDto, VideoDto>()
         .ForMember(dest => dest.VideoId, opt => opt.MapFrom(src => src.StreamerId))
         .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Name ?? "Live Stream"))
         .ForMember(dest => dest.Description, opt => opt.MapFrom(src => ""))
         .ForMember(dest => dest.Source, opt => opt.MapFrom(src => SystemPlatforms.Twitch.ToString()))
         .ForMember(dest => dest.Views, opt => opt.MapFrom(src => src.Viewers))
         .ForMember(dest => dest.VideoUploadStatus, opt => opt.MapFrom(src => VideoStatus.Uploading))
         .ForMember(dest => dest.IsPublished, opt => opt.MapFrom(src => true))
         .ForMember(dest => dest.Duration, opt => opt.MapFrom(src => 0))
         .ForMember(dest => dest.PreviewUrl, opt => opt.MapFrom(src => src.PreviewUrl))
         .ForMember(dest => dest.VideoUrl, opt => opt.MapFrom(src => src.StreamUrl))
         .ForMember(dest => dest.VideoTags, opt => opt.Ignore())
         .ForMember(dest => dest.IsPrivate, opt => opt.MapFrom(src => false))
         .ForMember(dest => dest.VideoTags, opt => opt.MapFrom(src => src.StreamTags))
         .ForMember(dest => dest.User, opt => opt.MapFrom(src => new UserVideoDto
         {
            UserId = src.StreamerId,
            AvatarUrl = src.StreamerInformation.AvatarUrl,
            IsFollowed = false,
            Username = src.StreamerInformation.Username
         }))

         .ForMember(dest => dest.CreatedAt, opt => opt.MapFrom(src => src.StartedAt))
         .ForMember(dest => dest.MediaType, opt => opt.MapFrom(src => MediaType.LiveStream));
      
      CreateMap<TwitchItem, StreamDto>()
         .ForMember(dest => dest.StreamerId, opt => opt.MapFrom(src => PlatformPrefixesConstants.TwitchPrefix + src.UserId))
         .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Title))
         .ForMember(dest => dest.StartedAt, opt => opt.MapFrom(src => src.StartedAt))
         .ForMember(dest => dest.Description, opt => opt.MapFrom(src => src.Description))
         .ForMember(dest => dest.Viewers, opt => opt.MapFrom(src => src.ViewerCount))
         .ForMember(dest => dest.Source, opt => opt.MapFrom(src => "Twitch"))
         .ForMember(dest => dest.StreamUrl, opt => opt.MapFrom(src => DefaultLinks.DefaultTwitchPlayerWatchLink + (src.UserName ?? src.BroadcasterLogin)))
         .ForMember(dest => dest.PreviewUrl, opt => opt.MapFrom(src => src.ThumbnailUrl))
         .ForMember(dest => dest.StreamTags, opt => opt.MapFrom(src => src.StreamTags))
         .ForMember(dest => dest.StreamerInformation, opt => opt.MapFrom(src => new UserStreamDto
         {
            Username = src.UserName ?? src.DisplayName ?? "Unknown",
            AvatarUrl = "" 
         }));
         
   }
}