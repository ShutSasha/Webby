using System.Net;
using Microsoft.AspNetCore.Mvc;
using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Response;
using Webby.AuthService.Interfaces.Services;
using Webby.AuthService.Models;

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
   public async Task<ActionResult<ApiResponse>> SignUp([FromBody] RegisterUserRequest request)
   {
      bool isRedirectionOnConfrimationPage = await _authService.Register(request);

      return isRedirectionOnConfrimationPage switch
      {
         false => Ok(ApiResponse.Ok("User successfully registered")),
         true => Ok(ApiResponse.Ok("Redirecting to confirmation email page")),
      };
   }

   [HttpPost("sign-in")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>>SignIn([FromBody] LoginUserRequest request)
   {
      var user = await _authService.Login(request);
      return Ok(ApiResponse<LoginUserResponse>.Ok("Login success", user));
   }

   [HttpPost("resend-verification-code")]
   public async Task<ActionResult> ResendVerificationCode([FromBody] ResendVerificationCodeRequest request)
   {
      await _authService.SendCode(request);
      return Ok(ApiResponse.Ok("Verification code sent"));
   }

   [HttpPost("verify-user")]
   public async Task<ActionResult<ApiResponse>> VerifyUser([FromBody] VerifyUserRequest request)
   {
      await _authService.VerifyEmail(request);
      return Ok(ApiResponse.Ok("User has been successfully verified"));
   }

   [HttpPost("google-auth")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>> PerformGoogleAuth([FromBody] GoogleAuthRequest request)
   {
      var user = await _authService.PerformGoogleAuth(request);
      return Ok(ApiResponse<LoginUserResponse>.Ok("Login success", user));
   }

   [HttpPost("refresh")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>> RefreshToken()
   {
      var token = HttpContext.Request.Headers["Authorization"].ToString().Split(" ")[1];
      var userRefreshResponse = await _authService.RefreshToken(token);
      return Ok(ApiResponse<LoginUserResponse>.Ok("Successfully refresh access token", userRefreshResponse));
   }

   [HttpGet("check")]
   public Task<ActionResult<ApiResponse>> Check()
   {
      return Task.FromResult<ActionResult<ApiResponse>>(Ok(ApiResponse.Ok("Server is running")));
   }
   
}