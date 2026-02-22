using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.AspNetCore.Mvc.Rendering;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Models;

namespace Webby.AuthService.Helpers.Jwt;

public class JwtProvider(IOptions<JwtOptions> options): IJwtProvider
{  
   private readonly JwtOptions _options = options.Value;  
  
   public string GenerateAccessToken(User user)
   {
      Claim[] claims = new[]
      {
         new Claim("Id", user.UserId.ToString()),
         new Claim("type", "access"),
         new Claim("Email", user.Email),

      }; 
      var signingCredentials = new SigningCredentials(  
         new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_options.AccessSecretKey)),  
         SecurityAlgorithms.HmacSha256);

      var accessToken = new JwtSecurityToken(
         claims: claims,
         signingCredentials: signingCredentials,
         expires: DateTime.Now.AddMinutes(_options.AccessExpiresDuration));
  
      var tokenValue = new JwtSecurityTokenHandler().WriteToken(accessToken);  
  
      return tokenValue;  
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