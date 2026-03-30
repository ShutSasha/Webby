using System.Text.Json.Serialization;
using Webby.NotificationService.Extensions;
using Webby.NotificationService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();

services.AddCorsPolicy("AllowApiGateway");
services.AddSwaggerConfig();
services.AddDbConnection(configuration);

services.AddRepositories();
services.AddServices();

services.AddControllers().AddJsonOptions(options =>
{
    options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
});

var app = builder.Build();

app.UseCors("AllowApiGetaway");

app.UseMiddleware<ExceptionMiddleware>();
app.UseMiddleware<ValidationExceptionMiddleware>();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger(c =>
    {
        c.RouteTemplate = "/docs/notification-service/{documentName}/swagger.json";
    });

    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/notification-service/v1/swagger.json", "Video Service API");
    });
}

app.UseRouting();
app.UseHttpsRedirection();
app.MapControllers();

app.Run();