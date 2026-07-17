using Microsoft.AspNetCore.Mvc;
using Microsoft.IdentityModel.JsonWebTokens;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Dtos.Complaint;
using Webby.UserService.Helpers.Jwt;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Controllers;

[ApiController]
[Route("/api/complaints")]
public class ComplaintController : ControllerBase
{
   private readonly IComplaintService _complaintService;
   
   public ComplaintController(IComplaintService complaintService)
   {
      _complaintService = complaintService;
   }

   [HttpGet("{userId:guid}")]
   [SwaggerOperation("Get user complaints","MODERATOR ROLE CLAIM REQUIRED")]
   public async Task<ActionResult<ApiResponse<List<ComplaintDto>>>> GetUserComplaints([FromRoute] Guid userId)
   {
      var userComplaints = await _complaintService.GetUserComplaints(userId);
      return Ok(ApiResponse<List<ComplaintDto>>.Ok("Successfully retrieved user complaints", userComplaints));
   }
   
   [HttpPost]
   [SwaggerOperation("Leave a complaint to user or video","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> CreateUserComplaints([FromBody] CreateComplaintRequest request)
   {
      var authorId = JwtHelper.ExtractUserId(HttpContext)!;
      await _complaintService.CreateComplaint(authorId.Value,request);
      return Ok(ApiResponse.Ok("Successfully create complaint"));
   }
   
}