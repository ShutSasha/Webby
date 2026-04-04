using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.AspNetCore.SignalR;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using Webby.NotificationService.Data;
using Webby.NotificationService.Helpers.Jwt;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Interfaces.Services;
using Webby.NotificationService.Providers;
using Webby.NotificationService.Repositories;

namespace Webby.NotificationService.Extensions;

public static class ApiExtension
{
   public static void AddCorsPolicy(this IServiceCollection serviceCollection, string policyName)
   {
      serviceCollection.AddCors(corsOptions =>
      {
         corsOptions.AddPolicy(policyName, policy =>
         {
            policy.WithHeaders().AllowCredentials();
            policy.WithHeaders().AllowAnyHeader();
            policy.WithOrigins("http://localhost:5000")
               .AllowAnyMethod()
               .AllowAnyHeader();
         });
      });
   }
   
   public static IServiceCollection AddSwaggerConfig(this IServiceCollection services)
   {
      services.AddSwaggerGen(options =>
      {
         
         options.EnableAnnotations();
         
         options.SwaggerDoc("v1", new OpenApiInfo
         {
            Title = "Notification Service API",
            Version = "v1"
         });
         
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

   public static void AddJwtAuthorization(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
         .AddJwtBearer(options =>
         {
            options.TokenValidationParameters = new TokenValidationParameters
            {
               ValidateIssuer = false,
               ValidateAudience = false,
               ValidateLifetime = true,
               ValidateIssuerSigningKey = true,
               IssuerSigningKey = new SymmetricSecurityKey(
                  Encoding.UTF8.GetBytes(configuration["JwtOptions:AccessSecretKey"]!)),
               ClockSkew = TimeSpan.Zero
            };

            options.Events = new JwtBearerEvents
            {
               OnMessageReceived = context =>
               {
                  var accessToken = context.Request.Query["accessToken"];
                  var path = context.HttpContext.Request.Path;

                  if (!string.IsNullOrEmpty(accessToken) && 
                      path.StartsWithSegments("/hubs/notifications"))
                  {
                     context.Token = accessToken;
                  }
                  return Task.CompletedTask;
               }
            };
         });
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
      serviceCollection.AddScoped<INotificationRepository, NotificationRepository>();
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<INotificationService, Services.NotificationService>();
   }

   public static void AddProviders(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddSingleton<IUserIdProvider, UserProvider>();
   }
   
}