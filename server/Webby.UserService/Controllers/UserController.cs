using Microsoft.AspNetCore.Mvc;
using Webby.UserService.Dtos;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Controllers;

[ApiController]
[Route("api/users")]
public class UserController : ControllerBase
{
   private readonly IUserService _userService;
   
   public UserController(IUserService userService)
   {
      _userService = userService;
   }

   [HttpGet("{userId:guid}")]
   public async Task<ActionResult<ApiResponse<UserDto>>> GetUserInformation([FromRoute] Guid userId)
   {
      var userExtractionResult = await _userService.GetUserInformation(userId);
      return Ok(ApiResponse.Ok("Successfully exctract user", userExtractionResult));
   }
   
}