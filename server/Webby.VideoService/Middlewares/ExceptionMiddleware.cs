using Webby.UserService.Helpers.Response;
using Webby.VideoService.Helpers.Exception;

namespace Webby.VideoService.Middlewares;

public class ExceptionMiddleware
{
   private readonly RequestDelegate _next;
   private readonly ILogger<ExceptionMiddleware> _logger;

   public ExceptionMiddleware(RequestDelegate next, ILogger<ExceptionMiddleware> logger)
   {
      _next = next;
      _logger = logger;
   }
   
   public async Task Invoke(HttpContext context)
   {
      try
      {
         await _next(context);
      }
      catch (ApiException ex)
      {
         var errorResponse = ApiResponse.Fail(ex.Message, ex.Errors);

         context.Response.StatusCode = ex.StatusCode;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
      catch (Exception)
      {
         var errorResponse = ApiResponse.Fail("Internal server error");

         context.Response.StatusCode = 500;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
   }
}