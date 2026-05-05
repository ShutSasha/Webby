using System.Reflection.Metadata.Ecma335;
using System.Text.Json.Serialization;
using Webby.UserService.Data.Configurations;

namespace Webby.UserService.Models;

public class Achievement
{
   public Guid AchievementId { get; set; }
   public string Title { get; set; }
   public string EventType { get; set; }
   public string Description { get; set; }
   public string IconUrl { get; set; }
   public int TargetValue { get; set; }
   
   [JsonIgnore]
   public ICollection<UserAchievement> UserAchievements { get; set; }
   
   [JsonIgnore]
   public ICollection<UserAchievementProgress> UserAchievementsProgresses { get; set; }
}