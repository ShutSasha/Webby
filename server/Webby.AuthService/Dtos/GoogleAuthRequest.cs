using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Dtos;

public class GoogleAuthRequest
{
   [Required]
   public Guid Id { get; set; }
   
   [Required]
   public string Name { get; set; }
   
   [Required]
   public string Email { get; set; }
   
   [Required]
   public string Image { get; set; }
}