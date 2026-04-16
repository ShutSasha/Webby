using System.ComponentModel;
using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;
using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.Stream.Enums;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Controllers;

[ApiController]
[Route("api/streams")]
public class StreamController : ControllerBase
{
   private readonly IStreamService _streamService;

   public StreamController(IStreamService streamService)
   {
      _streamService = streamService;
   }

   [HttpGet("search")]
   [SwaggerOperation("Search stream")]
   public async Task<ActionResult<ApiResponse<PagedResponse<StreamDto>>>> SearchStream([FromQuery] SearchStreamOptions searchStreamOptions)
   {
      var searchStreamResult = await _streamService.SearchStream(searchStreamOptions);
      return Ok(ApiResponse<PagedResponse<StreamDto>>.Ok("Successfully retrieved streams", searchStreamResult));
   }

   [HttpGet("{streamerId}")]
   [SwaggerOperation("Get stream information by streamer id")]
   public async Task<ActionResult<ApiResponse<StreamDto>>> GetStreamInformation([FromRoute] string streamerId,
      [FromQuery(Name = "searchPlatform")]
      SearchStreamPlatforms? searchStreamPlatforms)
   {
      var getStreamResult = await _streamService.GetStreamById(streamerId, searchStreamPlatforms);
      return Ok(ApiResponse<StreamDto>.Ok("Successfully retrieved stream", getStreamResult));
   }
   
}