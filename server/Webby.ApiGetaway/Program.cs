using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Ocelot.DependencyInjection;
using Ocelot.Middleware;
using Serilog;
using Webby.ApiGetaway.Extensions;
using Webby.ApiGetaway.Helpers.Jwt;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;
    builder.AddCustomSerilog();

    builder.WebHost.ConfigureKestrel(options => { options.Limits.MaxRequestBodySize = 5L * 1024 * 1024 * 1024; });
    
    services.AddHealthChecks();
    
    services.AddOpenApi();
    services.AddCorsPolicy("AllowWebOrigin");
    services.Configure<JwtOptions>(configuration.GetSection(nameof(JwtOptions)));

    services.AddJwtAuthentication(builder.Configuration);
    configuration.RegisterApiConfig(builder.Environment);

    services.AddOcelot(configuration);
    services.AddEndpointsApiExplorer();
    services.AddSwaggerInfo();

    var app = builder.Build();
    app.UseCustomSerilogRequestLogging();

    app.UseCors("AllowWebOrigin");
    app.UseAuthentication();
    app.UseAuthorization();
    
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        app.UseSwagger();
        app.UseSwaggerUI(c =>
        {
        c.SwaggerEndpoint("/auth/swagger/v1/swagger.json", "AuthService");
        c.SwaggerEndpoint("/room-category/swagger/v1/swagger.json", "RoomCategoryService");
        c.SwaggerEndpoint("/room-queue/swagger/v1/swagger.json", "RoomQueueService");
        c.SwaggerEndpoint("/vote/swagger/v1/swagger.json", "VoteService");
        c.SwaggerEndpoint("/user/swagger/v1/swagger.json", "UserService");
        c.SwaggerEndpoint("/video/swagger/v1/swagger.json", "VideoService");
        c.SwaggerEndpoint("/room/swagger/v1/swagger.json", "RoomService");
        c.SwaggerEndpoint("/chat/swagger/v1/swagger.json", "ChatService");
        c.SwaggerEndpoint("/ws/swagger/v1/swagger.json", "WsGateway");
        c.SwaggerEndpoint("/notification/swagger/v1/swagger.json", "NotificationService");
        c.SwaggerEndpoint("/achievement/swagger/v1/swagger.json", "AchievementService");
        c.RoutePrefix = "";
        });
    });
    app.UseApiExceptionHandling();

    app.UseWebSockets();
    
    app.MapHealthChecks("/livez", new HealthCheckOptions
    {
        Predicate = _ => false
    });

    app.MapHealthChecks("/readyz", new HealthCheckOptions
    {
        Predicate = _ => false
    });
    
    app.MapWhen(
        ctx => !ctx.Request.Path.StartsWithSegments("/livez") && 
               !ctx.Request.Path.StartsWithSegments("/readyz"),
        appBuilder => appBuilder.UseOcelot().Wait()
    );
    
    app.Run();
}
finally
{
    Log.CloseAndFlush();
}
