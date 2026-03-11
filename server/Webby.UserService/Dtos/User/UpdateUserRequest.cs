using System.ComponentModel.DataAnnotations;

namespace Webby.UserService.Dtos.User;

public class UpdateUserRequest
{
   [Required]
   public Guid UserId { get; set; }
   
   [MaxLength(150, ErrorMessage = "The maximum length can be no more than 150 characters")]
   public string? About { get; set; }
}