using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Services;

public interface ITokenService
{
   Task<AuthToken> GenerateToken(User user);
   Task<Guid> ExtractUserInfo(string accessToken);
   IEnumerable<string> ParseUserRolesFromToken(string accessToken);
}