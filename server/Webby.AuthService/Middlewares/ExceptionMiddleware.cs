using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Helpers.Response;

namespace Webby.AuthService.Middlewares;

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
         _logger.LogWarning("API Exception: {Message} (StatusCode: {StatusCode})", ex.Message, ex.StatusCode);

         var errorResponse = ApiResponse.Fail(ex.Message, ex.Errors);

         context.Response.StatusCode = ex.StatusCode;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
      catch (Exception ex)
      {
         _logger.LogError(ex, "Unhandled server error occurred");
         
         var errorMsg = new Dictionary<string, string> { { "errorMsg", ex.Message } };
         var errorResponse = ApiResponse.Fail("Internal server error", errorMsg);
         
         context.Response.StatusCode = 500;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
   }
}