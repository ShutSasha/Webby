namespace Webby.VideoService.Dtos.Video;

public class UploadVideoRequest
{
   public IFormFile VideoFile { get; set; } = null!;
}