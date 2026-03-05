namespace Webby.UserService.Models;

public class UserPremium
{
   public Guid UserId { get; set; }
   public DateTime ExpiresAt { get; set; }
}