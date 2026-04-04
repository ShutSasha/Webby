namespace Webby.VideoService.Dtos.User;

public class UserVideoDto
{
   public string UserId { get; set; }
   public required string Username { get; set; }
   public required string AvatarUrl { get; set; }
   public bool IsFollowed { get; set; }
}