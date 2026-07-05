using System.Security.Claims;
using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Helpers;

public interface IJwtProvider
{
   AuthToken GenerateAccessToken(User user);
   ClaimsPrincipal GetPrincipal(string accessToken);
   IEnumerable<string> GetUserRoles(string accessToken);
}