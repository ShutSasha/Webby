using Org.BouncyCastle.Asn1.IsisMtt.X509;
using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Services;

public interface IAuthService
{
   Task<bool> Register(RegisterUserRequest request);
   Task<LoginUserResponse> Login(LoginUserRequest request);
   Task SendCode(ResendVerificationCodeRequest request);
   Task VerifyEmail(VerifyUserRequest request);
   Task<LoginUserResponse> PerformGoogleAuth(GoogleAuthRequest request);
   Task<LoginUserResponse> RefreshToken(string accessToken);
   Task ChangeUserPassword(Guid userId, ChangeUserPasswordRequest request);
   Task ResetUserPassword(ResetUserPasswordRequest request);

}