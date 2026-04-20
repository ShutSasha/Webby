using System.Text.Json.Serialization;
using Amazon.S3;
using Grpc.Net.Client.Web;
using UserService;
using Webby.MediaService.GrpcServer;
using Webby.VideoService.Extensions;
using Webby.VideoService.Helpers.Seed;
using Webby.VideoService.Middlewares;
using Webby.VideoService.Services.Grpc;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

builder.WebHost.ConfigureKestrel(options =>
{
    options.Limits.MaxRequestBodySize = 5L * 1024 * 1024 * 1024; 
});

services.AddEndpointsApiExplorer();

services.AddCorsPolicy("AllowApiGetaway");

services.AddSwaggerConfig();
services.AddDbConnection(configuration);

services.AddAutoMapper(cfg =>
{
    cfg.LicenseKey = configuration["AutoMapper:LicenseKey"];
}, typeof(Program));

services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));

services.ConfigureOptionDependencies(configuration);

services.AddGrpc(options =>
{
    options.Interceptors.Add<GrpcExceptionInterceptor>();
});
services.AddRepositories();
services.AddServices();
services.AddExternalServices();

services.ConfigureGrpcConnections();

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

app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
app.MapGrpcService<VideoGrpcService>().EnableGrpcWeb();
app.MapGrpcService<MediaGrpcService>().EnableGrpcWeb();

app.UseRouting();
app.MapControllers();

app.Run();

