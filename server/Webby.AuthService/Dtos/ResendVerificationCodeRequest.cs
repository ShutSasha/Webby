using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class ResendVerificationCodeRequest
{
   [Required]
   [EmailAddress]
   public string Email { get; set; }
}