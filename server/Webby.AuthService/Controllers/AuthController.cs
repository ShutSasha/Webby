using System.IdentityModel.Tokens.Jwt;
using System.Net;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.DependencyModel.Resolution;
using Swashbuckle.AspNetCore.Annotations;
using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Helpers.Jwt;
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

   [HttpGet("is-valid-role")]
   [SwaggerOperation("Checks whether user has a valid role", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<bool>>> IsRoleValid()
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      var token = HttpContext.Request.Headers.Authorization.ToString().Substring("Bearer ".Length);
      var isValidRoleResult = await _authService.IsValidRole(userId, token);
      return Ok(ApiResponse<bool>.Ok("Successfully retrieved", isValidRoleResult));
   }
    
   [HttpPost("sign-up")]
   [SwaggerOperation("Registration route")]
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
   [SwaggerOperation("Login route")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>>SignIn([FromBody] LoginUserRequest request)
   {
      var user = await _authService.Login(request);
      return Ok(ApiResponse<LoginUserResponse>.Ok("Login success", user));
   }

   [HttpPost("resend-verification-code")]
   [SwaggerOperation("Resending verification code to user email")]
   public async Task<ActionResult> ResendVerificationCode([FromBody] ResendVerificationCodeRequest request)
   {
      await _authService.SendCode(request);
      return Ok(ApiResponse.Ok("Verification code sent"));
   }

   [HttpPost("verify-user")]
   [SwaggerOperation("User verification code check")]
   public async Task<ActionResult<ApiResponse>> VerifyUser([FromBody] VerifyUserRequest request)
   {
      await _authService.VerifyEmail(request);
      return Ok(ApiResponse.Ok("User has been successfully verified"));
   }

   [HttpPost("google-auth")]
   [SwaggerOperation("Proccessing google account data after OAuth verification step")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>> PerformGoogleAuth([FromBody] GoogleAuthRequest request)
   {
      var user = await _authService.PerformGoogleAuth(request);
      return Ok(ApiResponse<LoginUserResponse>.Ok("Login success", user));
   }
   
   [HttpPost("refresh")]
   [SwaggerOperation("Access token exchanging")]
   public async Task<ActionResult<ApiResponse<LoginUserResponse>>> RefreshToken()
   {
      var authHeader = HttpContext.Request.Headers.Authorization.ToString();

      if (string.IsNullOrWhiteSpace(authHeader))
         throw new ApiException("Refresh token error", 400, "Authorization header is missing.");

      if (!authHeader.StartsWith("Bearer ", StringComparison.OrdinalIgnoreCase))
         throw new ApiException("Refresh token error", 400, "Authorization header must use Bearer scheme.");

      var token = authHeader["Bearer ".Length..].Trim();

      if (string.IsNullOrWhiteSpace(token))
         throw new ApiException("Refresh token error", 400, "Token is missing.");

      var response = await _authService.RefreshToken(token);

      return Ok(ApiResponse<LoginUserResponse>
         .Ok("Successfully refreshed access token", response));
   }

   [HttpPatch("change-password")]
   [SwaggerOperation("Change user password route","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> ChangeUserPassword([FromBody] ChangeUserPasswordRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext);
      
      await _authService.ChangeUserPassword(userId,request);
      return Ok(ApiResponse.Ok("Successfully changed user password"));
   }
   
   [HttpPatch("reset-password")]
   [SwaggerOperation("Reset user password after user email verification")]
   public async Task<ActionResult<ApiResponse>> ResetUserPassword([FromBody] ResetUserPasswordRequest request)
   {
      await _authService.ResetUserPassword(request);
      return Ok(ApiResponse.Ok("Successfully reset user password"));
   }
   
}