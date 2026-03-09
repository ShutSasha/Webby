using Webby.UserService.Dtos;
using Webby.UserService.Dtos.User;

namespace Webby.UserService.Interfaces.Service;

public interface IUserService
{
   Task<UserProfileResponse> GetUserInformation(Guid userId);
   Task<UserDto> UpdateUserInformation(UpdateUserRequest request);
   Task<UserDto> EditUserIcon(Guid id, string fileName, Stream fileStream, string contentType);
   Task FollowUser(UserFollowRequest request);
   Task UnfollowUser(UserFollowRequest request);
   Task UnlockAchievement(Guid userId, Guid achievementId);
   Task PinUserAchievement(Guid userId, Guid achievementId);
   Task UnpinUserAchievement(Guid userId, Guid achievementId);
}