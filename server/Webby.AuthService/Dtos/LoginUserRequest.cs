using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class LoginUserRequest
{
   [Required]
   public string Email { get; set; }
   
   [Required]
   public string Password { get; set; }
}