namespace Webby.AchievementService.Dtos.Achievement;

public class AchievementDto : ProfileAchievementDto
{
   public string Description { get; set; }
   public int AchievementProgressValue { get; set; }
   public int TargetValue { get; set; }
   public bool IsUnlocked { get; set; }
}