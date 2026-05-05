namespace Webby.UserService.Models;

public class UserAchievementProgress
{
   public Guid UserId { get; init; }
   public User User { get; set; }
   public Guid AchievementId { get; init; }
   public int CurrentValue { get; set; }
   public Achievement Achievement { get; set; }
}