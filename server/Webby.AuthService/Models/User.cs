namespace Webby.AuthService.Models;

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
   public Guid? PremiumId { get; set; }
}