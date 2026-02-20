namespace Webby.AuthService.Dtos;

public class UserDto
{
   public Guid UserId { get; set; }
   public string Email { get; set; }
   public string Username { get; set; }
   public string About { get; set; }
   public string AvatarUrl { get; set; }
}