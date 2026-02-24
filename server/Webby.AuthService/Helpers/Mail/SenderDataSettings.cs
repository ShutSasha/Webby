using System.ComponentModel.DataAnnotations;

namespace Webby.AuthService.Helpers.Mail;

public class SenderDataSettings
{
   public string SenderEmail { get; set; }
   public string SenderPassword { get; set; }
}