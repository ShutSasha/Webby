namespace Webby.AchievementService.Models;

public class UserAchievementProgress
{
   public Guid UserId { get; init; }
   public Guid AchievementId { get; init; }
   public int CurrentValue { get; set; }
   public Achievement Achievement { get; set; }
}