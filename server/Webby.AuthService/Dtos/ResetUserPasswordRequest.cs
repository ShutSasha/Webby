namespace Webby.AuthService.Dtos;

public class ResetUserPasswordRequest
{
   public string Email { get; set; }
   public string NewPassword { get; set; }
}