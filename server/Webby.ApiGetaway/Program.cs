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

    if (app.Environment.IsDevelopment())
    {
        app.UseSwagger();
        app.UseSwaggerUI(c =>
        {
        c.SwaggerEndpoint("/auth/swagger/v1/swagger.json", "AuthService");
        c.SwaggerEndpoint("/room-category/swagger/v1/swagger.json", "RoomCategoryService");
        c.SwaggerEndpoint("/user/swagger/v1/swagger.json", "UserService");
        c.SwaggerEndpoint("/video/swagger/v1/swagger.json", "VideoService");
        c.SwaggerEndpoint("/room/swagger/v1/swagger.json", "RoomService");
        c.SwaggerEndpoint("/room-queue/swagger/v1/swagger.json", "RoomQueueService");
        c.SwaggerEndpoint("/chat/swagger/v1/swagger.json", "ChatService");
        c.SwaggerEndpoint("/notification/swagger/v1/swagger.json", "NotificationService");
        c.RoutePrefix = "";
        });
    }
    app.UseApiExceptionHandling();

    app.UseWebSockets();
    await app.UseOcelot();
    app.Run();
}
finally
{
    Log.CloseAndFlush();
}
