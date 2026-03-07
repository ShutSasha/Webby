using Webby.UserService.Dtos.User;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IUserRepository : IRepository<User>
{
   Task<UserFollowStats> GetUserFollowBlock(Guid userId);
}