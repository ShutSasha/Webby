using System.Text.Json.Serialization;
using Webby.NotificationService.Extensions;
using Webby.NotificationService.Hubs;
using Webby.NotificationService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();

services.AddAuthorization();

services.AddCorsPolicy("AllowApiGateway");
services.AddSwaggerConfig();
services.AddDbConnection(configuration);

builder.Services.AddSignalR(options =>
{
    options.KeepAliveInterval = TimeSpan.FromSeconds(10);
    options.ClientTimeoutInterval = TimeSpan.FromSeconds(30);
});

services.AddRepositories();
services.AddServices();
services.AddProviders();

services.AddControllers().AddJsonOptions(options =>
{
    options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
});

var app = builder.Build();

app.UseMiddleware<ExceptionMiddleware>();
app.UseMiddleware<ValidationExceptionMiddleware>();

app.UseHttpsRedirection();
app.UseRouting();
app.UseCors("AllowApiGateway");
app.UseAuthentication();
app.UseAuthorization();


if (app.Environment.IsDevelopment())
{
    app.UseSwagger(c =>
    {
        c.RouteTemplate = "/docs/notification-service/{documentName}/swagger.json";
    });

    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/notification-service/v1/swagger.json", "Notification Service API");
    });
}

app.MapControllers();
app.MapHub<NotificationHub>("/hubs/notifications");

app.Run();