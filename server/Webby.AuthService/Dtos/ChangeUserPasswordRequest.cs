using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Mvc.Infrastructure;

namespace Webby.AuthService.Dtos;

public class ChangeUserPasswordRequest
{
   
   [Required]
   public string NewPassword { get; set; }
   
   [Required]
   public string CurrentPassword { get; set; }
}