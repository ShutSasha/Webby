namespace Webby.AuthService.Dtos;

public class AuthToken
{
   public string AccessToken { get; set; }
   public long ExpiresAt { get; set; }
}