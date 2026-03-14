namespace Webby.UserService.Dtos.Achievement;

public class GetUserAchievementsResponse
{
   public List<AchievementDto> PinnedAchievementDtos { get; set; }
   public List<AchievementDto> Achievements { get; set; }
}