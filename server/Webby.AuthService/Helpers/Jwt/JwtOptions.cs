namespace Webby.AuthService.Helpers.Jwt;

public class JwtOptions  
{  
   public string AccessSecretKey { get; set; } = string.Empty;  
   public int AccessExpiresDuration { get; set; } = 30;
}