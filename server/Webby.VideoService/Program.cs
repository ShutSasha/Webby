using System.Text.Json.Serialization;
using Amazon.S3;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Serilog;
using Webby.VideoService.Data;
using Webby.VideoService.Extensions;
using Webby.VideoService.Middlewares;
using Webby.VideoService.Services.Grpc;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;
    builder.AddCustomSerilog();
    
    services.AddHealthChecks()
        .AddNpgSql(configuration.GetConnectionString(nameof(AppDbContext)), tags: new[] { "ready" })
        .AddRedis(configuration.GetConnectionString("Redis"), tags: new[] { "ready" });
    
    builder.WebHost.ConfigureKestrel(options => { options.Limits.MaxRequestBodySize = 5L * 1024 * 1024 * 1024; });

    services.AddEndpointsApiExplorer();

    services.AddCorsPolicy("AllowApiGetaway");

    services.AddSwaggerConfig();
    services.AddDbConnection(configuration);
    services.ConfigureRedisConnection(configuration);

    services.AddAutoMapper(cfg => { cfg.LicenseKey = configuration["AutoMapper:LicenseKey"]; }, typeof(Program));

    services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));

    services.ConfigureOptionDependencies(configuration);
    
    services.AddHelpers();
    services.AddGrpc(options => { options.Interceptors.Add<GrpcExceptionInterceptor>(); });
    services.AddRepositories();
    services.AddServices();
    services.AddExternalServices();
    services.AddBackgroundServices();
    services.AddInterceptors();

    services.ConfigureGrpcConnections(configuration);

    services.AddControllers().AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
    });

    var app = builder.Build();
    app.UseCustomSerilogRequestLogging();
    app.UseCors("AllowApiGetaway");

    app.UseMiddleware<ExceptionMiddleware>();
    app.UseMiddleware<ValidationExceptionMiddleware>();
    
    app.UseSwagger(c => { c.RouteTemplate = "docs/video-service/{documentName}/swagger.json"; });
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/video-service/v1/swagger.json", "Video Service API");
        c.RoutePrefix = "docs/video-service";
    });

    app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
    app.MapGrpcService<VideoGrpcService>().EnableGrpcWeb();
    app.MapGrpcService<MediaGrpcService>().EnableGrpcWeb();

    app.UseRouting();
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
