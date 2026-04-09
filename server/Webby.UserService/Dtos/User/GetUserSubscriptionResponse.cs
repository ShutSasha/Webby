
namespace Webby.UserService.Dtos.User;

public class GetUserSubscriptionResponse
{
   public DateTime ExpiresAt { get; set; }
   public bool IsSubscriptionExpired => ExpiresAt < DateTime.UtcNow;
}