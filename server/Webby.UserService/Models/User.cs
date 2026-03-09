using Webby.UserService.Models.Enums;

namespace Webby.UserService.Models;

public class User
{
   public Guid UserId { get; set; }
   public string Email { get; set; }
   public string Username { get; set; }
   public string? Password { get; set; }
   public string About { get; set; }
   public string AvatarUrl { get; set; }
   public string VerificationCode { get; set; }
   public bool isVerified { get; set; }
   public Role Role { get; set; }
   
   public ICollection<Complaint> Complaints { get; set; }
   public ICollection<UserFollower> Followers { get; set; }
   public ICollection<UserAchievement> UserAchievements { get; set; }
   public ICollection<UserFollower> Following { get; set; }
}