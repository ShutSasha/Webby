using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Middlewares;

public class ValidationExceptionMiddleware
{
   private readonly RequestDelegate _next;

   public ValidationExceptionMiddleware(RequestDelegate next)
   {
      _next = next;
   }

   public async Task Invoke(HttpContext context)
   {
      if (context.Request.ContentType?.StartsWith("application/grpc") == true)
      {
         await _next(context);
         return;
      }

      var originalBody = context.Response.Body;

      await using var memStream = new MemoryStream();
      context.Response.Body = memStream;

      try
      {
         await _next(context);

         memStream.Seek(0, SeekOrigin.Begin);

         if (context.Response.StatusCode == 400 &&
             memStream.Length > 0 &&
             context.Response.ContentType?.Contains("application/problem+json") == true)
         {
            var json = await new StreamReader(memStream).ReadToEndAsync();
            var details = JsonSerializer.Deserialize<ValidationProblemDetails>(json);

            var errors = details.Errors.ToDictionary(
               e => e.Key.ToLower(),
               e => e.Value.First()
            );

            var formatted = ApiResponse.Fail("Validation error", errors);
            var output = JsonSerializer.Serialize(formatted, new JsonSerializerOptions
            {
               PropertyNamingPolicy = JsonNamingPolicy.CamelCase
            });

            context.Response.ContentType = "application/json";
            context.Response.ContentLength = output.Length;
            context.Response.Body = originalBody;

            await context.Response.WriteAsync(output);
            return;
         }

         memStream.Seek(0, SeekOrigin.Begin);
         await memStream.CopyToAsync(originalBody);
      }
      finally
      {
         context.Response.Body = originalBody;
      }
   }
}