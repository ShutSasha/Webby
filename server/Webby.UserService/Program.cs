using System.Text.Json.Serialization;
using Webby.UserService.Extensions;
using Webby.UserService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();
services.AddSwaggerGen();

services.AddCorsPolicy("AllowApiGetaway");
services.AddDbConnection(configuration);

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
        c.RouteTemplate = "docs/user-service/{documentName}/swagger.json";
    });

    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/user-service/v1/swagger.json", "User Service API");
    });
}

app.UseHttpsRedirection();
app.UseRouting();

app.MapControllers();
app.Run();