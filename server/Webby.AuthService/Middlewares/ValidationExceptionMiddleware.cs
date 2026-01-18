using System.Text.Json;
using Microsoft.AspNetCore.Mvc;

public class ValidationExceptionMiddleware
{
   private readonly RequestDelegate _next;

   public ValidationExceptionMiddleware(RequestDelegate next)
   {
      _next = next;
   }

   public async Task Invoke(HttpContext context)
   {
      var originalBody = context.Response.Body;

      using var memStream = new MemoryStream();
      context.Response.Body = memStream;

      await _next(context);

      if (context.Response.StatusCode == 400 &&
          memStream.Length > 0 &&
          context.Response.ContentType?.Contains("application/problem+json") == true)
      {
         memStream.Seek(0, SeekOrigin.Begin);
         var json = await new StreamReader(memStream).ReadToEndAsync();

         var details = JsonSerializer.Deserialize<ValidationProblemDetails>(json);

         var errors = details.Errors.ToDictionary(
            e => e.Key,
            e => e.Value.First()
         );

         var formatted = new
         {
            message = "Validation failed",
            status = 400,
            errors
         };

         var output = JsonSerializer.Serialize(formatted);

         context.Response.ContentType = "application/json";
         context.Response.Body = originalBody;

         await context.Response.WriteAsync(output);
         return;
      }

      memStream.Seek(0, SeekOrigin.Begin);
      await memStream.CopyToAsync(originalBody);
   }
}