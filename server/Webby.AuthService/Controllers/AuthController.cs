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
      bool isRedirectionOnConfrimationPage = await _authService.Register(request);

      return isRedirectionOnConfrimationPage switch
      {
         false => Ok(ApiResponse.Ok("User successfully registered")),
         true => Ok(ApiResponse.Ok("Redirecting to confirmation email page")),
      };
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

   [HttpPost("google-auth")]
   public async Task<IActionResult> PerformGoogleAuth([FromBody] GoogleAuthRequest request)
   {
      var user = await _authService.PerformGoogleAuth(request);
      return Ok(ApiResponse.Ok("Login success", user));
   }

   [HttpGet("check")]
   public async Task<IActionResult> Check()
   {
      return Ok(ApiResponse.Ok("Server is running"));
   }
   
}