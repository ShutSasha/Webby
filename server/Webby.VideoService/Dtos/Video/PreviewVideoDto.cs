using Webby.VideoService.Dtos.User;
using Webby.VideoService.Interfaces.Dto;

namespace Webby.VideoService.Dtos.Video;

public class PreviewVideoDto: IVideoDtoWithUser
{
   public string VideoId { get; init; }
   public string Name { get; init; }
   public string PreviewUrl { get; init; }
   public int Views { get; init; }
   public bool IsPrivate { get; init; }
   public DateTime CreatedAt { get; init; }
   public long Duration { get; init; }
   public UserVideoDto User { get; set; }
}