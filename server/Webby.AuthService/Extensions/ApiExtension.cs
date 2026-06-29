using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using Webby.AuthService.Data;
using Webby.AuthService.Helpers;
using Webby.AuthService.Helpers.Jwt;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Interfaces.Repositories;
using Webby.AuthService.Interfaces.Services;
using Webby.AuthService.Models;
using Webby.AuthService.Repositories;
using Webby.AuthService.Services;

namespace Webby.AuthService.Extensions;

public static class ApiExtension
{
   public static void AddCorsPolicy(this IServiceCollection serviceCollection, string policyName)
   {
      serviceCollection.AddCors(corsOptions =>
      {
         corsOptions.AddPolicy(policyName, policy =>
         {
            policy
               .WithOrigins("http://localhost:5000")
               .AllowAnyMethod()
               .AllowCredentials()
               .AllowAnyHeader();
         });
      });
   }
   
    public static IServiceCollection AddSwaggerConfig(this IServiceCollection services)
    {
       services.AddSwaggerGen(options =>
       {
          options.EnableAnnotations();
 
          options.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
          {
             Name = "Authorization",
             Type = SecuritySchemeType.Http,
             Scheme = "Bearer",
             BearerFormat = "JWT",
             In = ParameterLocation.Header,
             Description = "Enter ONLY your JWT token"
          });
 
          options.AddSecurityRequirement(new OpenApiSecurityRequirement
          {
             {
                new OpenApiSecurityScheme
                {
                   Reference = new OpenApiReference
                   {
                      Type = ReferenceType.SecurityScheme,
                      Id = "Bearer"
                   }
                },
                Array.Empty<string>()
             }
          });
       });
 
       return services;
    }

   public static void AddDbConnection(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddDbContext<AppDbContext>(options =>
      {
         options.UseNpgsql(configuration.GetConnectionString(nameof(AppDbContext)));
      });
   }

   public static void AddRepositories(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IRepository<User>, GenericRepository<User>>();
      serviceCollection.AddScoped<IAuthRepository, AuthRepository>();
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IMailService, MailService>();
      serviceCollection.AddScoped<IAuthService, Services.AuthService>();
      serviceCollection.AddScoped<ITokenService, TokenService>();
   }

   public static void AddHelpers(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IPasswordHasher, PasswordHasher>();
      serviceCollection.AddScoped<IJwtProvider, JwtProvider>();
      
   }
   
}