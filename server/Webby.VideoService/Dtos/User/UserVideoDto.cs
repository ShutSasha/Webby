namespace Webby.VideoService.Dtos.User;

public class UserVideoDto : UserVideoDtoBase
{
   public required string AvatarUrl { get; set; }
   public bool IsFollowed { get; set; }
}