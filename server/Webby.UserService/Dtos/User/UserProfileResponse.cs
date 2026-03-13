using Webby.UserService.Dtos.Achievement;

namespace Webby.UserService.Dtos.User;

public class UserProfileResponse
{
   public UserDto User { get; set; }
   public UserFollowStats UserFollowStats { get; set; }
   public List<AchievementDto> PinnedUserAchievements { get; set; }

}