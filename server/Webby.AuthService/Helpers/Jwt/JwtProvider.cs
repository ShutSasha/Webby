using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.AspNetCore.Mvc.Rendering;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using Webby.AuthService.Dtos;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Models;

namespace Webby.AuthService.Helpers.Jwt;

public class JwtProvider(IOptions<JwtOptions> options): IJwtProvider
{  
   private readonly JwtOptions _options = options.Value;  
  
   public AuthToken GenerateAccessToken(User user)
   {
      var expiration = DateTime.UtcNow.AddMinutes(_options.AccessExpiresDuration);
      var roles = new List<string>();

      switch(user.Role)
      {
         case Role.Admin:
            roles.Add("Admin");
            roles.Add("Moderator");
            break;
         case Role.Moderator:
            roles.Add("Moderator");
            break;
         default:
            roles.Add("User");
            break;
      }
      
      var claims = new List<Claim>
      {
         new Claim("Id", user.UserId.ToString())
      };
      
      claims.AddRange(roles.Select(r => new Claim("Role", r)));

      var signingCredentials = new SigningCredentials(
         new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_options.AccessSecretKey)),
         SecurityAlgorithms.HmacSha256);
      
      var token = new JwtSecurityToken(
         claims: claims,
         signingCredentials: signingCredentials,
         expires: expiration
      );
      
      return new AuthToken
      {
         AccessToken = new JwtSecurityTokenHandler().WriteToken(token),
         ExpiresAt = new DateTimeOffset(expiration).ToUnixTimeSeconds()
      };
   }
  
   public ClaimsPrincipal GetPrincipal(string accessToken)  
   {  
      var securityKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_options.AccessSecretKey));  
        
      var validation = new TokenValidationParameters  
      {  
         IssuerSigningKey = securityKey,  
         ValidateIssuer = false,  
         ValidateAudience = false,  
         ValidateLifetime = false,  
         ValidateIssuerSigningKey = true  
      };  
  
      return new JwtSecurityTokenHandler().ValidateToken(accessToken, validation, out _);  
   }  
  
}