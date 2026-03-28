using System.ComponentModel.DataAnnotations;
using Webby.UserService.Models.Enums;

namespace Webby.UserService.Dtos.Complaint;

public class CreateComplaintRequest
{
   
   [Required]
   public Guid TargetId { get; set; }
   
   [Required]
   public required string ReasonType { get; set; }

   [Required]
   [EnumDataType(typeof(ComplaintTargetType), ErrorMessage = "Selected target type is invalid")]
   public ComplaintTargetType? TargetType { get; set; }
   public string? AdditionalInfo { get; set; }
}