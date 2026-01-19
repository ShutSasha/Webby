using Microsoft.AspNetCore.Mvc;
using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Response;
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
      return Ok(ApiResponse.Ok("User successfully registered"));
   }

   [HttpPost("sign-in")]
   public async Task<IActionResult> SignIn([FromBody] LoginUserRequest request)
   {
      var user = await _authService.Login(request);
      return Ok(ApiResponse.Ok("Login success", user));
   }

   [HttpPost("resend-verification-code")]
   public async Task<IActionResult> ResendVerificationCode([FromBody] ResendVerificationCodeRequest request)
   {
      await _authService.SendCode(request);
      return Ok(ApiResponse.Ok("Verification code sent"));
   }

   [HttpPost("verify-user")]
   public async Task<IActionResult> VerifyUser([FromBody] VerifyUserRequest request)
   {
      var user = await _authService.VerifyEmail(request);
      return Ok(ApiResponse.Ok("User has been successfully verified", user));
   }
   
}