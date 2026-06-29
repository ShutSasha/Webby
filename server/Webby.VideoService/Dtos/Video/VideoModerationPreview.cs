using Webby.VideoService.Dtos.User;
using Webby.VideoService.Interfaces.Dto;

namespace Webby.VideoService.Dtos.Video;

public class VideoModerationPreview : IVideoDtoWithUser
{
   public string VideoId { get; set; }
   public string PreviewUrl { get; set; }
   public string Name { get; set; }
   public bool IsBanned { get; set; }
   public DateTime CreatedAt { get; set; }
   public UserVideoDto? User { get; set; }
}