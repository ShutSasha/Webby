using System.Text.Json.Serialization;
using Amazon.S3;
using Serilog;
using Webby.UserService.Extensions;
using Webby.UserService.Middlewares;
using UserGrpcService = Webby.UserService.Services.Grpc.UserGrpcService;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;
    builder.AddCustomSerilog();

    services.AddEndpointsApiExplorer();
    services.AddSwaggerGen();

    services.AddCorsPolicy("AllowApiGetaway");
    services.AddSwaggerConfig();
    services.AddDbConnection(configuration);
    services.ConfigureRedisConnection(configuration);


    services.AddAutoMapper(cfg => { cfg.LicenseKey = configuration["AutoMapper:LicenseKey"]; }, typeof(Program));

    services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));

    services.ConfigureOptionDependencies(configuration);

    services.ConfigureGrpcConnections();

    services.AddBackgroundWorkers();
    services.AddRepositories();
    services.AddServices();
    services.AddHelpers();
    services.AddGrpc();


    services.AddControllers().AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter());
    });

    var app = builder.Build();

    app.UseCustomSerilogRequestLogging();

    app.UseCors("AllowApiGetaway");

    app.UseMiddleware<ExceptionMiddleware>();
    app.UseMiddleware<ValidationExceptionMiddleware>();

    if (app.Environment.IsDevelopment())
    {
        app.UseSwagger(c => { c.RouteTemplate = "docs/user-service/{documentName}/swagger.json"; });

        app.UseSwaggerUI(c => { c.SwaggerEndpoint("/docs/user-service/v1/swagger.json", "User Service API"); });
    }

    app.UseRouting();

    app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
    app.MapGrpcService<UserGrpcService>().EnableGrpcWeb();


    app.MapControllers();

    app.Run();
}
finally
{
    Log.CloseAndFlush();
}