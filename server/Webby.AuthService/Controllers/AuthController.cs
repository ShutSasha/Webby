using Microsoft.AspNetCore.Mvc;
using Webby.AuthService.Dtos;
using Webby.AuthService.Interfaces.Services;

namespace Webby.AuthService.Controllers;

[ApiController]
[Route("api/auth")]
public class AuthController : ControllerBase
{
   private readonly IAuthService _authService;
   
   public AuthController(IAuthService authService)
   {
      _authService = authService;
   }

   [HttpPost("sign-up")]
   public async Task<IActionResult> SignUp([FromBody] RegisterUserRequest request)
   {
      await _authService.Register(request);

      return Ok(new { Message = "User successfully registered" });
      
   }
   
   [HttpPost("sign-in")]
   public async Task<IActionResult> SignIn([FromBody] LoginUserRequest request)
   {
      var userLoginResult = await _authService.Login(request);

      return Ok(userLoginResult);
      
   }

   [HttpPost("resend-verification-code")]
   public async Task<IActionResult> ResendVerificationCode([FromBody] ResendVerificationCodeRequest request)
   {
      await _authService.SendCode(request);
      return Ok(new { Message = "Verification code successfully send" });
   }

   [HttpPost("verify-user")]
   public async Task<IActionResult> VerifyUser([FromBody] VerifyUserRequest request)
   {
      var userVerifyResult = await _authService.VerifyEmail(request);
      return Ok(userVerifyResult);
   }
   
}