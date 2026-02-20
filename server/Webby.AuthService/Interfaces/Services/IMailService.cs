namespace Webby.AuthService.Interfaces.Services;

public interface IMailService
{
   Task SendVerificationCode(string recieverEmail, string code);
   
}