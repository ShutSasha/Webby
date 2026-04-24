using Webby.UserService.Dtos;
using Webby.UserService.Dtos.Search;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Service;

public interface IUserService
{
   Task<UserProfileResponse> GetUserInformation(Guid userId);
   Task<UserDto> UpdateUserInformation(UpdateUserRequest request);
   Task<UserDto> EditUserIcon(Guid id, string fileName, Stream fileStream, string contentType);
   Task<string> ProcessFollow(UserFollowRequest request);
   Task UnlockAchievement(Guid userId, Guid achievementId);
   Task PinUserAchievement(Guid userId, Guid achievementId);
   Task UnpinUserAchievement(Guid userId, Guid achievementId);
   Task<List<UserFollowersDto>> GetUserFollowers(Guid userId);
   Task<List<UserFollowersDto>> GetUserFollows(Guid userId);
   Task<User?> GetById(Guid userId);
   Task<UserFollowingResponse> IsUserFollowing(Guid userId, Guid targetId);
   Task<GetUserSubscriptionResponse?> GetUserSubscription(Guid userId, bool isOwner);
   Task<PagedResponse<UserDto>> SearchUsers(SearchOptions searchOptions);
   

}