namespace Webby.AchievementService.Dtos.Achievement;

public class AchievementDto : ProfileAchievementDto
{
   public string Description { get; set; }
   public bool IsUnlocked { get; set; }
}