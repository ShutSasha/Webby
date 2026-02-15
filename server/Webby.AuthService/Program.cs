using Webby.AuthService.Extensions;
using Webby.AuthService.Helpers.Mail;
using Webby.AuthService.Middlewares;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
var configuration = builder.Configuration;

services.AddEndpointsApiExplorer();
services.AddSwaggerGen();

services.AddCorsPolicy("AllowApiGetaway");
services.AddDbConnection(configuration);

services.AddControllers();
services.Configure<SenderDataSettings>(configuration.GetSection("SenderData"));

services.AddRepositories();
services.AddHelpers();
services.AddServices();
services.AddAutoMapper(AppDomain.CurrentDomain.GetAssemblies());
services.AddOpenApi();

var app = builder.Build();

app.UseCors("AllowApiGetaway");
app.UseMiddleware<ExceptionMiddleware>();
app.UseMiddleware<ValidationExceptionMiddleware>();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseRouting();

app.MapControllers();
app.Run();
