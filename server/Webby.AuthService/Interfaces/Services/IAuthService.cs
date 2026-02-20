using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Services;

public interface IAuthService
{
   Task<bool> Register(RegisterUserRequest request);
   Task<UserDto> Login(LoginUserRequest request);
   Task SendCode(ResendVerificationCodeRequest request);
   Task<UserDto> VerifyEmail(VerifyUserRequest request);
   Task<UserDto> PerformGoogleAuth(GoogleAuthRequest request);
}