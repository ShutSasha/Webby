namespace Webby.UserService.Dtos.User;

public class UserProfileResponse
{
   public UserDto User { get; set; }
   public UserFollowStats UserFollowStats { get; set; }
}