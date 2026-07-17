using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Interfaces.Services;
using Webby.AuthService.Models;

namespace Webby.AuthService.Services;

public class TokenService: ITokenService
{
   private readonly IJwtProvider _jwtProvider;
   
   public TokenService(IJwtProvider jwtProvider)
   {
      _jwtProvider = jwtProvider;
   }

   public Task<AuthToken> GenerateToken(User user) 
      => Task.FromResult(_jwtProvider.GenerateAccessToken(user));

   public Task<Guid> ExtractUserInfo(string accessToken)
   {
      var errors = new Dictionary<string, string>();
      var principal = _jwtProvider.GetPrincipal(accessToken);

      if (principal == null || principal.Claims.All(c => c.Type != "Id"))
      {
         errors["token"] = "Invalid access token";
         throw new ApiException("Extract user error", 400, errors);
      }

      var userId = principal.Claims.FirstOrDefault(c => c.Type == "Id")?.Value;

      if (string.IsNullOrEmpty(userId))
      {
         errors["userId"] = "Invalid token payload parameters";
         throw new ApiException("Extract user error",400,errors);
      }

      return Task.FromResult(Guid.Parse(userId));
   }

   public IEnumerable<string> ParseUserRolesFromToken(string accessToken)
      => _jwtProvider.GetUserRoles(accessToken);
}