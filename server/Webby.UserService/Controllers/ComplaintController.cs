using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Dtos.Complaint;
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
   [SwaggerOperation("Leave a complaint to user","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> CreateUserComplaints([FromBody] CreateUserComplaintRequest request)
   {
      await _complaintService.CreateUserComplaint(request);
      return Ok(ApiResponse.Ok("Successfully create user complaints"));
   }
   
}