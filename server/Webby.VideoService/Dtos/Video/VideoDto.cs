using Webby.VideoService.Dtos.User;
using Webby.VideoService.Models.Enums;
using Webby.VideoService.Services.Background;

namespace Webby.VideoService.Dtos.Video;

public class VideoDto
{
   public Guid VideoId { get; set; }
   public required string Name { get; set; }
   public string? Description { get; set; }
   public int Views { get; set; }
   public required string PreviewUrl { get; set; }
   public required bool IsPrivate { get; set; }
   public DateTime CreatedAt { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   public TimeSpan Duration { get; set; }
   public List<string>? VideoTags { get; set; }
   public UserVideoDto? User { get; set; }
}