using Serilog;
using Serilog.Sinks.SystemConsole.Themes;

namespace Webby.AchievementService.Extensions;

public static class LoggerExtension
{
   public static WebApplicationBuilder AddCustomSerilog(this WebApplicationBuilder builder)
   {
      builder.Logging.ClearProviders();

      builder.Host.UseSerilog((context, services, configuration) => configuration
         .ReadFrom.Configuration(context.Configuration)
         .ReadFrom.Services(services)
         .Enrich.FromLogContext()
         .WriteTo.Console(
            outputTemplate: "[{Timestamp:HH:mm:ss} {Level:u3}] [ACHIEVEMENT-SERVICE] {Message:lj}{NewLine}{Exception}",
            theme: AnsiConsoleTheme.Sixteen
         ));

      return builder;
   }
   
   public static WebApplication UseCustomSerilogRequestLogging(this WebApplication app)
   {
      app.UseSerilogRequestLogging(options =>
      {
         options.MessageTemplate = 
            @"{{
  ""method"": ""{RequestMethod}"",
  ""path"": ""{RequestPath}"",
  ""status"": {StatusCode},
  ""duration_ms"": {Elapsed:0.0000},
  ""trace_id"": ""{TraceId}""
}}";

         options.EnrichDiagnosticContext = (diagnosticContext, httpContext) =>
         {
            var traceId = System.Diagnostics.Activity.Current?.TraceId.ToString() 
                          ?? httpContext.TraceIdentifier;
                          
            diagnosticContext.Set("TraceId", traceId);
         };
      });

      return app;
   }
}