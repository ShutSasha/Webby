namespace Webby.UserService.Dtos.User;

public class UserFollowRequest
{
   public Guid UserId { get; set; }
   public Guid FollowerId { get; set; }
}