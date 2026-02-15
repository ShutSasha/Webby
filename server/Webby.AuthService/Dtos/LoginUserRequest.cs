using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;
public class LoginUserRequest
{
   [Required]
   public string Email { get; set; }
   
   public string Password { get; set; }

}