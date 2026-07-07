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
            UserId = src.StreamerInformation.UserId,
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
            UserId = src.UserId,
            Username = src.UserName ?? src.DisplayName ?? "Unknown",
            AvatarUrl = "" 
         }));
      
      CreateMap<YouTubeVideoResponse.Item, VideoDto>()
         .ForMember(dest => dest.VideoId, opt => opt.MapFrom((src, dest) => PlatformPrefixesConstants.YouTubePrefix + src.Id.ToString()))
         .ForMember(dest => dest.Name, opt => opt.MapFrom(src => src.Snippet.Title))
         .ForMember(dest => dest.Description, opt => opt.MapFrom(src => src.Snippet.Description))
         .ForMember(dest => dest.Source, opt => opt.MapFrom(src => "YouTube"))
         .ForMember(dest => dest.VideoUrl, opt => opt.MapFrom((src, dest) => DefaultLinks.BaseWatchLinkUrl + src.Id.ToString()))
         .ForMember(dest => dest.CreatedAt, opt => opt.MapFrom(src => src.Snippet.PublishedAt))
         .ForMember(dest => dest.MediaType, opt => opt.MapFrom(src => MediaType.Video))
         .ForMember(dest => dest.IsPrivate, opt => opt.MapFrom(src => false))
         .ForMember(dest => dest.VideoUploadStatus, opt => opt.MapFrom(src => VideoStatus.Ready))
         .ForMember(dest => dest.VideoTags, opt => opt.MapFrom(src => src.Snippet.Tags))
         .ForMember(dest => dest.PreviewUrl, opt => opt.MapFrom(src => 
            src.Snippet.Thumbnails.Medium != null ? src.Snippet.Thumbnails.Medium.Url : 
            src.Snippet.Thumbnails.Default != null ? src.Snippet.Thumbnails.Default.Url : ""))         
         .ForMember(dest => dest.Duration, opt => opt.MapFrom(src => 
            !string.IsNullOrEmpty(src.ContentDetails.Duration) 
               ? (long)Math.Round(System.Xml.XmlConvert.ToTimeSpan(src.ContentDetails.Duration).TotalSeconds) 
               : 0L))
         .ForMember(dest => dest.Views, opt => opt.MapFrom((src, dest) => 
            src.Statistics.ViewCount != null && int.TryParse(src.Statistics.ViewCount.ToString(), out var v) ? v : 0))
         .ForMember(dest => dest.User, opt => opt.MapFrom(src => new UserVideoDto
         {
            UserId = src.Snippet.ChannelId ?? "",
            Username = src.Snippet.ChannelTitle,
            AvatarUrl = "",
            IsFollowed = false
         }));

      CreateMap<Models.Video, VideoModerationPreview>();
      
   }
}