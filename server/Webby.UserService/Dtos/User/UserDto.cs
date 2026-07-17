using Webby.UserService.Models.Enums;

namespace Webby.UserService.Dtos.User;

public class UserDto
{
   public Guid UserId { get; set; }
   public string Email { get; set; }
   public string Username { get; set; }
   public string About { get; set; }
   public string AvatarUrl { get; set; }
   public bool IsBanned { get; set; }
   public DateTime CreatedAt { get; set; }
   public Role Role { get; set; }
}