using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Repositories;

public interface IAuthRepository : IRepository<User>
{
   Task<bool> HasUserRoles(Guid userId, IEnumerable<string> tokenRoles);
}