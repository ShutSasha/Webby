namespace Webby.AuthService.Dtos;

public class ChangeUserPasswordRequest
{
   public Guid UserId { get; set; }
   public string NewPassword { get; set; }
}