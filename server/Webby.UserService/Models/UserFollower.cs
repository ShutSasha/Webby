namespace Webby.UserService.Models;

public class UserFollower
{
   public Guid UserId { get; set; }
   public Guid FollowerId { get; set; }
   public User FollowedUser { get; set; }
   public User FollowerUser { get; set; }
}