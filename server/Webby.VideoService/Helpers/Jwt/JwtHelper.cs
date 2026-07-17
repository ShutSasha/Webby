using System.IdentityModel.Tokens.Jwt;
using Webby.VideoService.Helpers.Exception;

namespace Webby.VideoService.Helpers.Jwt;

public static class JwtHelper
{
   public static Guid? ExtractUserId(HttpContext context, bool shouldThrowException = true)
   {
      var authHeader = context.Request.Headers["Authorization"].ToString();

      if (string.IsNullOrEmpty(authHeader) || !authHeader.StartsWith("Bearer "))
      {
         if (shouldThrowException)
            throw new ApiException("Extract user error", 400, "Auth header missing");

         return null;
      }

      var token = authHeader.Substring("Bearer ".Length);

      try
      {
         var handler = new JwtSecurityTokenHandler();
         var jwtToken = handler.ReadJwtToken(token);

         var claim = jwtToken.Claims.FirstOrDefault(c => c.Type == "Id");

         if (claim == null || !Guid.TryParse(claim.Value, out var userId))
         {
            if (shouldThrowException)
               throw new ApiException("Extract user error", 400, "Invalid token");

            return null;
         }

         return userId;
      }
      catch
      {
         if (shouldThrowException)
            throw new ApiException("Extract user error", 400, "Invalid token");

         return null;
      }
   }
}