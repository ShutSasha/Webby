namespace Webby.UserService.Dtos.User;

public class UserFollowersDto
{
   public Guid UserId { get; set; }
   public string AvatarUrl { get; set; }
   public string Username { get; set; }
   public int FollowersCount { get; set; }
}