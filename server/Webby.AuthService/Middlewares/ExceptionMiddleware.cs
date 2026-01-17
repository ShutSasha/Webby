using Webby.AuthService.Helpers.Exception;

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
         _logger.LogError(ex, "API Exception: {Message}", ex.Message);

         var errorResponse = new
         {
            status = ex.StatusCode,
            message = ex.Message
         };

         context.Response.StatusCode = ex.StatusCode;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
      catch (Exception ex)
      {
         _logger.LogError(ex, "Unhandled Exception: {Message}", ex.Message);

         var errorResponse = new
         {
            status = 500,
            message = "An unexpected error occurred."
         };

         context.Response.StatusCode = 500;
         context.Response.ContentType = "application/json";

         await context.Response.WriteAsJsonAsync(errorResponse);
      }
   }
}