using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class RegisterUserRequest
{
   [Required]
   [EmailAddress]
   public string Email { get; set; }
   
   [Required]
   [StringLength(25, MinimumLength = 3, ErrorMessage = "{0} must be between {2} and {1} characters length.")]
   public string Username { get; set; }
   
   [Required]
   [StringLength(25, MinimumLength = 3, ErrorMessage = "{0} must be between {2} and {1} characters length.")]
   public string Password { get; set; }
}