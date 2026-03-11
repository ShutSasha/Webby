using System.ComponentModel.DataAnnotations;

namespace Webby.UserService.Dtos.Complaint;

public class CreateUserComplaintRequest
{
   
   [Required]
   public Guid AuthorId { get; set; }
   
   [Required]
   public Guid TargetUserId { get; set; }
   
   [Required]
   public string ReasonType { get; set; }

   public string? AdditionalInfo { get; set; }
}