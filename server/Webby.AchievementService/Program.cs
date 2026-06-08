using System.Text.Json.Serialization;
using Amazon.S3;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Serilog;
using Webby.AchievementService.Data;
using Webby.AchievementService.Extensions;
using Webby.AchievementService.Middlewares;
using AchievementGrpcService = Webby.AchievementService.Services.Grpc.AchievementGrpcService;

try
{
   var builder = WebApplication.CreateBuilder(args);
   var services = builder.Services;
   var configuration = builder.Configuration;
   builder.AddCustomSerilog();

   services.AddHealthChecks()
      .AddNpgSql(configuration.GetConnectionString(nameof(AppDbContext)), tags: new[] { "ready" })
      .AddRedis(configuration.GetConnectionString("Redis"), tags: new[] { "ready" });
      
   services.AddEndpointsApiExplorer();
   services.AddSwaggerGen();
   
   services.ConfigureOptionDependencies(configuration);
   
   services.AddCorsPolicy("AllowApiGetaway");
   services.AddSwaggerConfig();
   services.AddDbConnection(configuration);
   services.ConfigureRedisConnection(configuration);
   
   services.AddAutoMapper(cfg => { cfg.LicenseKey = configuration["AutoMapper:LicenseKey"]; }, typeof(Program));
   services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));
   
   services.AddControllers().AddJsonOptions(options =>
   {
      options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
   });
   
   services.AddHelpers();
   services.AddBackgroundWorkers();
   services.AddInterceptors();
   services.AddRepositories();
   services.AddServices();
   
   services.ConfigureGrpcConnection();

   
   var app = builder.Build();
   
   app.UseCustomSerilogRequestLogging();

   app.UseCors("AllowApiGetaway");
   
   app.UseMiddleware<ExceptionMiddleware>();
   app.UseMiddleware<ValidationExceptionMiddleware>();

   if (app.Environment.IsDevelopment())
   {
      app.UseSwagger(c => { c.RouteTemplate = "docs/achievement-service/{documentName}/swagger.json"; });

      app.UseSwaggerUI(c => { c.SwaggerEndpoint("/docs/achievement-service/v1/swagger.json", "Achievement Service API"); });
   }
   
   app.UseRouting();
   
   app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
   app.MapGrpcService<AchievementGrpcService>();
   app.MapControllers();
   
   app.MapHealthChecks("/livez", new HealthCheckOptions
   {
      Predicate = _ => false
   });

   app.MapHealthChecks("/readyz", new HealthCheckOptions
   {
      Predicate = check => check.Tags.Contains("ready")
   });

   app.Run(); 
}
finally
{
   Log.CloseAndFlush();
}
