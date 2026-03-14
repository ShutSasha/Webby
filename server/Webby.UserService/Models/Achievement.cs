using System.Text.Json.Serialization;
using Webby.UserService.Data.Configurations;

namespace Webby.UserService.Models;

public class Achievement
{
   public Guid AchievementId { get; set; }
   public string Code { get; set; }
   public string Title { get; set; }
   public string Description { get; set; }
   public string IconUrl { get; set; }
   
   [JsonIgnore]
   public ICollection<UserAchievement> UserAchievements { get; set; }
}