using System.IdentityModel.Tokens.Jwt;
using Webby.VideoService.Helpers.Exception;

namespace Webby.VideoService.Helpers.Jwt;

public static class JwtHelper
{
   public static Guid ExtractUserId(HttpContext context)
   {
      var authHeader = context.Request.Headers["Authorization"].ToString();

      if (string.IsNullOrEmpty(authHeader) || !authHeader.StartsWith("Bearer "))
      {
         throw new ApiException("Extract user error",400,"Auth header missing");
      }

      var token = authHeader.Substring("Bearer ".Length);

      var handler = new JwtSecurityTokenHandler();
      var jwtToken = handler.ReadJwtToken(token);

      var claim = jwtToken.Claims.FirstOrDefault(c => c.Type == "Id");

      if (claim == null || !Guid.TryParse(claim.Value, out var userId))
      {
         throw new ApiException("Extract user error",400,"Invalid token");
      }

      return userId;
   }
}