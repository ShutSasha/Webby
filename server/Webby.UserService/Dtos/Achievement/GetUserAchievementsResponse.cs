namespace Webby.UserService.Dtos.Achievement;

public class GetUserAchievementsResponse
{
   public List<AchievementDto> PinnedAchievements { get; set; }
   public List<AchievementDto> Achievements { get; set; }
}