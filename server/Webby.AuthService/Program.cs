using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Serilog;
using Webby.AuthService.Data;
using Webby.AuthService.Extensions;
using Webby.AuthService.Helpers.Jwt;
using Webby.AuthService.Helpers.Mail;
using Webby.AuthService.Middlewares;

try
{
    var builder = WebApplication.CreateBuilder(args);
    var services = builder.Services;
    var configuration = builder.Configuration;

    builder.AddCustomSerilog();
    services.AddHealthChecks()
        .AddNpgSql(configuration.GetConnectionString(nameof(AppDbContext)), tags: new[] { "ready" });
    
    
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
    services.AddAutoMapper(cfg => { cfg.LicenseKey = configuration["AutoMapper:LicenseKey"]; }, typeof(Program));

    services.AddOpenApi();

    var app = builder.Build();

    app.UseCors("AllowApiGetaway");
    app.UseCustomSerilogRequestLogging();
    app.UseMiddleware<ExceptionMiddleware>();
    app.UseMiddleware<ValidationExceptionMiddleware>();
    
    app.UseSwagger(c => { c.RouteTemplate = "docs/auth-service/{documentName}/swagger.json"; });
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/docs/auth-service/v1/swagger.json", "Auth Service API");
        c.RoutePrefix = "docs/auth-service";
    });
    
    app.UseHttpsRedirection();
    app.UseRouting();

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
