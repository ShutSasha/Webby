using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Dtos.Video;

public class UploadVideoResponse
{
   public string VideoId { get; set; }
   public Guid UserId { get; set; }
   public required string Name { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   public bool IsPrivate { get; set; }
}