using AutoMapper;
using Webby.UserService.Dtos.Achievement;

namespace Webby.UserService.Dtos.User;

public class UserProfileResponse
{
   public UserDto User { get; set; }
   public UserFollowStats UserFollowStats { get; set; }
   public List<ProfileAchievementDto> PinnedUserAchievements { get; set; }
   public bool IsPremiumUser { get; set; }

}