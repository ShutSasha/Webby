using System.IdentityModel.Tokens.Jwt;
using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.NotificationService.Dtos.Notification;
using Webby.NotificationService.Dtos.Pagination;
using Webby.NotificationService.Helpers.Jwt;
using Webby.NotificationService.Helpers.Response;
using Webby.NotificationService.Interfaces.Services;
using Webby.NotificationService.Models;

namespace Webby.NotificationService.Controllers;

[ApiController]
[Route("api/notifications")]
public class NotificationController : ControllerBase
{
   private readonly INotificationService _notificationService;

   public NotificationController(INotificationService notificationService)
   {
      _notificationService = notificationService;
   }

   [HttpGet]
   [SwaggerOperation("Get all user notifications", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PagedResponse<Notification>>>> GetUserNotifications([FromQuery] PaginationRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var userNotifications = await _notificationService.GetUserNotifications(userId.Value, request);
      return Ok(ApiResponse<PagedResponse<Notification>>.Ok("Successfully retrieved user notifications",userNotifications));
   }
   
   [HttpGet("unread")]
   [SwaggerOperation("Get users unread notification","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PagedResponse<Notification>>>> GetUserUnreadNotification([FromQuery] PaginationRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var notifications = await _notificationService.GetUnreadNotifications(userId.Value, request);
      return Ok(ApiResponse<PagedResponse<Notification>>.Ok("Successfully retrieved user notifications", notifications));
   } 
   
   [HttpGet("read")]
   [SwaggerOperation("Get users read notification","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<PagedResponse<Notification>>>> GetUserReadNotification([FromQuery] PaginationRequest request)
   {
      var userId = JwtHelper.ExtractUserId(HttpContext)!;
      var notifications = await _notificationService.GetReadNotifications(userId.Value,request);
      return Ok(ApiResponse<PagedResponse<Notification>>.Ok("Successfully retrieved user notifications", notifications));
   }

   [HttpGet("count-unread-messages")]
   [SwaggerOperation("Get user unread messages count")]
   public async Task<ActionResult<ApiResponse<int>>> GetUnreadMessagesCount()
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext,shouldThrowException: false);
      var countOfUnreadMessages = await _notificationService.GetNotificationsCount(requestUserId);
      return Ok(ApiResponse<int>.Ok("Successfully retrieved count of unread messages", countOfUnreadMessages));
   }

   [HttpGet("count-notifications")]
   [SwaggerOperation("Count user unread and read messages", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<GetNotificationsCountResponse>>> GetNotificationsCount()
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext)!;
      var getNotificationsCountResponse = await _notificationService.GetUsersNotificationsCount(requestUserId.Value);
      return Ok(ApiResponse<GetNotificationsCountResponse>.Ok("Successfully retrieved user notifications",
         getNotificationsCountResponse));
   }
   
   [HttpPost]
   [SwaggerOperation("Create notification", "ADMIN ROLE REQUIRED")]
   public async Task<ActionResult<ApiResponse<Notification>>> CreateNotification([FromBody] CreateNotificationRequest request)
   {
      var notification = await _notificationService.CreateNotification(request);
      return Ok(ApiResponse<Notification>.Ok("Successfully create notification", notification));
   }

   [HttpPost("set-read-status")]
   [SwaggerOperation("Change notification read-status", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> ChangeNotificationReadStatus([FromBody] List<Guid> notificationIds)
   {
      var requestedUserId = JwtHelper.ExtractUserId(HttpContext)!;
      await _notificationService.ChangeReadStatus(requestedUserId.Value,notificationIds);
      return Ok(ApiResponse.Ok("Successfully change notifications read status"));
   }
   
   [HttpPatch]
   [SwaggerOperation("Update notification", "ADMIN ROLE REQUIRED")]
   public async Task<ActionResult<ApiResponse<Notification>>> UpdateNotification([FromBody] UpdateNotificationRequest request)
   {
      var notification = await _notificationService.UpdateNotification(request);
      return Ok(ApiResponse<Notification>.Ok("Successfully update notification", notification));
   }

   [HttpDelete("{notificationId:guid}")]
   [SwaggerOperation("Delete notification route", "AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse>> DeleteNotification([FromRoute] Guid notificationId)
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext)!;
      await _notificationService.DeleteNotification(requestUserId.Value, notificationId);
      return Ok(ApiResponse.Ok("Successfully delete video"));
   }
}