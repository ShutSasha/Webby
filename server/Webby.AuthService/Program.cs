using System.Text.Json.Serialization;
using Webby.AuthService.Extensions;
using Webby.AuthService.Helpers.Jwt;
using Webby.AuthService.Helpers.Mail;
using Webby.AuthService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();
services.AddSwaggerGen();
services.AddSwaggerConfig();

services.AddCorsPolicy("AllowApiGetaway");
services.AddDbConnection(configuration);

services.AddControllers().AddJsonOptions(options =>
{
    options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
});
services.Configure<SenderDataSettings>(configuration.GetSection("SenderData"));
services.Configure<JwtOptions>(configuration.GetSection(nameof(JwtOptions)));

services.AddRepositories();
services.AddHelpers();
services.AddServices();
services.AddAutoMapper(cfg =>
{
    cfg.LicenseKey = configuration["AutoMapper:LicenseKey"];
}, typeof(Program));

services.AddOpenApi();

var app = builder.Build();

app.UseCors("AllowApiGetaway");
app.UseMiddleware<ExceptionMiddleware>();
app.UseMiddleware<ValidationExceptionMiddleware>();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger(c =>
    {
        c.RouteTemplate = "docs/auth-service/{documentName}/swagger.json";
    });

    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/auth-service/v1/swagger.json", "Auth Service API");
    });
}

app.UseHttpsRedirection();
app.UseRouting();

app.MapControllers();
app.Run();
