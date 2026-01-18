using Microsoft.OpenApi.Models;
using Ocelot.DependencyInjection;
using Ocelot.Middleware;
using Webby.ApiGetaway.Extensions;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddOpenApi();
services.AddCorsPolicy("AllowWebOrigin");

configuration
    .SetBasePath(builder.Environment.ContentRootPath)
    .AddJsonFile("ocelot.json", optional: false, reloadOnChange: true)
    .AddEnvironmentVariables();

services.AddOcelot(configuration);
services.AddEndpointsApiExplorer();

services.AddSwaggerInfo();

var app = builder.Build();

app.UseCors("AllowWebOrigin");

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("http://localhost:5001/swagger/v1/swagger.json", "AuthService");
        c.RoutePrefix = "";
    });
}

await app.UseOcelot();
app.Run();
