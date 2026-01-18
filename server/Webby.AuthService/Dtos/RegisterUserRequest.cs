using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class RegisterUserRequest
{
   [Required]
   [EmailAddress]
   public string Email { get; set; }
   
   [Required]
   [Range(3, 25, ErrorMessage = "{0} must be between {1} and {2} characters length.")]
   public string Username { get; set; }
   
   [Required]
   [Range(3, 25, ErrorMessage = "{0} must be between {1} and {2} characters length.")]
   public string Password { get; set; }
}