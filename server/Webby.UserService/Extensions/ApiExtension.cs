using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using StackExchange.Redis;
using UserService.AchievementGrpcClient;
using Webby.UserService.Clients;
using Webby.UserService.Data;
using Webby.UserService.Dtos.Storage;
using Webby.UserService.Helpers.Notification;
using Webby.UserService.Helpers.Payment;
using Webby.UserService.Interfaces.Helpers;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Middlewares;
using Webby.UserService.Repositories;
using Webby.UserService.Services;

namespace Webby.UserService.Extensions;

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

   public static void AddInterceptors(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddTransient<GrpcExceptionInterceptor>();
      serviceCollection.AddTransient<GrpcClientExceptionInterceptor>();
   }

   public static void AddRepositories(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IUserRepository, UserRepository>();
      serviceCollection.AddScoped<IComplaintRepository, ComplaintRepository>();
      serviceCollection.AddScoped<IPaymentRepository, PaymentRepository>();
      serviceCollection.AddScoped<IUserPremiumRepository, UserPremiumRepository>();
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IStorageService,StorageService>();
      serviceCollection.AddScoped<IUserService,Services.UserService>();
      serviceCollection.AddScoped<IComplaintService, ComplaintService>();
      serviceCollection.AddScoped<IPaymentService, PaymentService>();
   }
   
   public static void AddHelpers(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<INotificationFactory,NotificationFactory>();
   }

   public static void ConfigureOptionDependencies(this IServiceCollection serviceCollection, IConfiguration config)
   {
      serviceCollection.Configure<AwsOptions>(config.GetSection(nameof(AwsOptions)));
      serviceCollection.Configure<PaymentSettings>(config.GetSection("PaymentOptions"));
   }
   
   public static void ConfigureGrpcConnections(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddGrpcClient<VideoGrpcService.VideoGrpcServiceClient>(o =>
      {
         o.Address = new Uri(configuration["GrpcClients:VideoServiceUrl"]);
      });
      serviceCollection.AddGrpcClient<NotificationService.GrpcClient.NotificationGrpcService.NotificationGrpcServiceClient>(o =>
      {
         o.Address = new Uri(configuration["GrpcClients:NotificationServiceUrl"]);
      });
      serviceCollection.AddGrpcClient<AchievementGrpcService.AchievementGrpcServiceClient>(o =>
      {
         o.Address = new Uri("http://localhost:5009");
      }).AddInterceptor<GrpcClientExceptionInterceptor>();
      
   }
   
   public static void ConfigureRedisConnection(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddSingleton<IConnectionMultiplexer>(sp =>
         ConnectionMultiplexer.Connect(configuration.GetConnectionString("Redis") ?? string.Empty));
   }
   
}