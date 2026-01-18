using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Services;

public interface IAuthService
{
   Task Register(RegisterUserRequest request);
   Task<User> Login(LoginUserRequest request);
   Task SendCode(ResendVerificationCodeRequest request);
   Task<User> VerifyEmail(VerifyUserRequest request);
   
}