using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.AchievementService.Helpers.Response;
using Webby.AchievementService.Interfaces.Services;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Controllers;

[ApiController]
[Route("api/events")]
public class EventTypeController : ControllerBase
{
   private readonly IEventTypeService _eventTypeService;
   
   public EventTypeController(IEventTypeService eventTypeService)
   {
      _eventTypeService = eventTypeService;
   }

   [HttpGet]
   [SwaggerOperation("Get accessible events on platform","MODERATION ROLE REQUIRED")]
   public async Task<ActionResult<ApiResponse<List<EventType>>>> GetEventTypes()
   {
      var eventTypes = await _eventTypeService.GetEventTypes();
      return Ok(ApiResponse<List<EventType>>.Ok("Successfully retrieved event types", eventTypes));
   }
}