namespace Webby.AuthService.Dtos;

public class LoginUserResponse
{
   public UserDto User { get; set; }
   public string AccessToken { get; set; }
   public long AccessTokenExpiresAt { get; set; }
}