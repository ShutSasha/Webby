using Webby.UserService.Dtos.User;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IUserRepository : IRepository<User>
{
   Task<UserFollowStats> GetUserFollowBlock(Guid userId);
   Task AddUserFollowing(Guid userId, Guid followerId);
   Task DeleteUserFollowing(Guid userId, Guid followerId);
   Task<bool> HasUserFollow(Guid userId, Guid followerId);
   Task<List<UserFollower>> GetUserFollows(Guid userId);
   Task<List<UserFollower>> GetUserFollowers(Guid userId);
   Task<List<User>> GetByIds(List<Guid> ids);
   Task<List<Guid>> GetSubscriptionIds(Guid requestUserId);
   Task<List<Guid>> SearchByUsername(List<Guid> userIds, string search, int limit, int skip);
   Task<bool> AreMutualFollowers(Guid user1Id, Guid user2Id);
}