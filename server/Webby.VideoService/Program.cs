using System.Text.Json.Serialization;
using Amazon.S3;
using Microsoft.OpenApi.Models;
using Webby.VideoService.Extensions;
using Webby.VideoService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();

services.AddCorsPolicy("AllowApiGetaway");
services.AddSwaggerConfig();
services.AddDbConnection(configuration);

services.AddAutoMapper(AppDomain.CurrentDomain.GetAssemblies());
services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));

services.ConfigureOptionDependencies(configuration);

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
        c.RouteTemplate = "docs/video-service/{documentName}/swagger.json";
    });

    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/video-service/v1/swagger.json", "Video Service API");
    });
}

app.UseHttpsRedirection();
app.UseRouting();

app.MapControllers();

app.Run();

