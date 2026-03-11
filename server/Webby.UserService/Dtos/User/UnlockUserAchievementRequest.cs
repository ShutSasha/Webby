using System.ComponentModel.DataAnnotations;

namespace Webby.UserService.Dtos.User;

public class UnlockUserAchievementRequest
{
   [Required] 
   public Guid UserId { get; set; }
   
   [Required] 
   public Guid AchievementId { get; set; }
}