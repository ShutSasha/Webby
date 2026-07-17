using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace Webby.UserService.Models;

public class UserPremium
{
   public Guid UserId { get; set; }
   public User User { get; set; }
   public DateTime ExpiresAt { get; set; }
}