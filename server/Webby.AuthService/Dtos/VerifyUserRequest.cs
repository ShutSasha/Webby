namespace Webby.AuthService.Dtos;

public class VerifyUserRequest
{
   public string Email { get; set; }
   public string VerificationCode { get; set; }
}