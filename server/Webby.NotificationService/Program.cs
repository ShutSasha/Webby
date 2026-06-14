using System.Text;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.IdentityModel.Tokens;
using Serilog;
using Webby.NotificationService.Data;
using Webby.NotificationService.Extensions;
using Webby.NotificationService.GrpcService;
using Webby.NotificationService.Hubs;
using Webby.NotificationService.Middlewares;
using Webby.NotificationService.Services;
using NotificationGrpcService = Webby.NotificationService.Services.Grpc.NotificationGrpcService;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;

    builder.AddCustomSerilog();

    services.AddHealthChecks()
        .AddNpgSql(configuration.GetConnectionString(nameof(AppDbContext)), tags: new[] { "ready" });
    
    services.AddAuthorization();
    services.AddMemoryCache();
    services.AddCorsPolicy("AllowApiGateway");
    services.AddDbConnection(configuration);
    services.AddSwaggerConfig();

    builder.Services.AddSignalR(options =>
    {
        options.KeepAliveInterval = TimeSpan.FromSeconds(10);
        options.ClientTimeoutInterval = TimeSpan.FromSeconds(30);
    }).AddJsonProtocol(options => { options.PayloadSerializerOptions.Converters.Add(new JsonStringEnumConverter()); });

    services.AddGrpc(options => { options.Interceptors.Add<GrpcExceptionInterceptor>(); });
    services.ConfigureGrpcConnections(configuration);
    
    services.AddRepositories();
    services.AddServices();
    services.AddProviders();
    services.AddJwtAuthorization(configuration);

    services.AddControllers().AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
    });


    var app = builder.Build();

    app.UseCustomSerilogRequestLogging();
    app.UseMiddleware<ExceptionMiddleware>();
    app.UseMiddleware<ValidationExceptionMiddleware>();

    app.UseHttpsRedirection();
    app.UseRouting();
    app.UseCors("AllowApiGateway");
    app.UseAuthentication();
    app.UseAuthorization();
    
    app.UseSwagger(c => { c.RouteTemplate = "/docs/notification-service/{documentName}/swagger.json"; });
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/notification-service/v1/swagger.json", "Notification Service API");
        c.RoutePrefix = "docs/notification-service";
    });

    app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
    app.MapGrpcService<NotificationGrpcService>().EnableGrpcWeb();
    app.MapControllers();
    app.MapHub<NotificationHub>("/hubs/notifications");

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