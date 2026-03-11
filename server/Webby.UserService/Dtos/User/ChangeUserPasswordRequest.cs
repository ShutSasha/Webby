namespace Webby.UserService.Dtos.User;

public class ChangeUserPasswordRequest
{
   public Guid UserId { get; set; }
   public string Password { get; set; }
}