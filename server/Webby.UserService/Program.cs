using System.Text.Json.Serialization;
using Amazon.S3;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Serilog;
using Webby.UserService.Data;
using Webby.UserService.Extensions;
using Webby.UserService.Middlewares;
using UserGrpcService = Webby.UserService.Services.Grpc.UserGrpcService;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;
    builder.AddCustomSerilog();
    
    services.AddHealthChecks()
        .AddNpgSql(configuration.GetConnectionString(nameof(AppDbContext)), tags: new[] { "ready" })
        .AddRedis(configuration.GetConnectionString("Redis"), tags: new[] { "ready" });
    
    services.AddEndpointsApiExplorer();
    services.AddSwaggerGen();

    services.AddCorsPolicy("AllowApiGetaway");
    services.AddSwaggerConfig();
    services.AddDbConnection(configuration);
    services.ConfigureRedisConnection(configuration);


    services.AddAutoMapper(cfg => { cfg.LicenseKey = configuration["AutoMapper:LicenseKey"]; }, typeof(Program));

    services.AddSingleton<IAmazonS3>(AwsS3ClientFactory.CreateS3Client(configuration));

    services.ConfigureOptionDependencies(configuration);

    
    services.AddInterceptors();
    services.ConfigureGrpcConnections(configuration);

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
    
    app.UseSwagger(c => { c.RouteTemplate = "docs/user-service/{documentName}/swagger.json"; });
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/user-service/v1/swagger.json", "User Service API");
        c.RoutePrefix = "docs/user-service";
    });

    app.UseRouting();

    app.UseGrpcWeb(new GrpcWebOptions { DefaultEnabled = true });
    app.MapGrpcService<UserGrpcService>().EnableGrpcWeb();


    app.MapControllers();
        
    app.MapHealthChecks("/livez", new HealthCheckOptions
    {
        Predicate = _ => false
    });

    app.MapHealthChecks("/readyz", new HealthCheckOptions
    {
        Predicate = check => check.Tags.Contains("ready")
    });

    app.Run();
}
finally
{
    Log.CloseAndFlush();
}