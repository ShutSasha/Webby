using System.ComponentModel.DataAnnotations;

namespace Webby.UserService.Dtos.Achievement;

public class UpdateAchievementRequest
{
   [Required] 
   public Guid AchievementId { get; set; }
   
   [Required] 
   public IFormFile File { get; set; }

   [Required]
   public string Title { get; set; }

   [Required]
   public string Description { get; set; }

   [Required]
   public string Code { get; set; }
}