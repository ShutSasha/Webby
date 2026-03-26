using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Models;

public class Video
{
   public Guid VideoId { get; set; }
   public Guid UserId { get; set; }
   public required string Name { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   public string? Description { get; set; }
   public int Views { get; set; }
   public DateTime CreatedAt { get; set; }
   public string? VideoUrl { get; set; }
   public string? PreviewUrl { get; set; }
   public bool IsPrivate { get; set; }
   public ICollection<VideoTag>? VideoTags { get; set; }
   
}