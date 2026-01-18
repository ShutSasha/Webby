using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class VerifyUserRequest
{
   [Required]
   public string Email { get; set; }
   
   [Required]
   public string VerificationCode { get; set; }
}