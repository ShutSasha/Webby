using System.Security.Claims;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Helpers;

public interface IJwtProvider
{
   string GenerateAccessToken(User user);
   ClaimsPrincipal GetPrincipal(string accessToken);
}