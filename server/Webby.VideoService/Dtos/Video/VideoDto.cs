using System.Text.Json.Serialization;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Interfaces.Dto;
using Webby.VideoService.Models.Enums;
using Webby.VideoService.Services.Background;

namespace Webby.VideoService.Dtos.Video;

public class VideoDto : IVideoDtoWithUser
{
   public string VideoId { get; set; }
   public required string Name { get; set; }
   public string? Description { get; set; }
   
   public string Source { get; set; } = SearchVideoPlatforms.Webby.ToString();
   public int Views { get; set; }
   public required string PreviewUrl { get; set; }
   public required bool IsPrivate { get; set; }
   public DateTime CreatedAt { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   public bool IsPublished { get; set; } = false;
   public long Duration { get; set; }
   public string? VideoUrl { get; set; }
   public List<string>? VideoTags { get; set; }
   public UserVideoDto? User { get; set; }
}