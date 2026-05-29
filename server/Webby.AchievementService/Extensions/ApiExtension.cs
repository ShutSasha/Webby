using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using StackExchange.Redis;
using Webby.AchievementService.Data;
using Webby.AchievementService.Dtos.Storage;
using Webby.AchievementService.Helpers.Notification;
using Webby.AchievementService.Interfaces.Helpers.Notification;
using Webby.AchievementService.Interfaces.Repositories;
using Webby.AchievementService.Interfaces.Services;
using Webby.AchievementService.Middlewares;
using Webby.AchievementService.Repositories;
using Webby.AchievementService.Services;
using Webby.AchievementService.Services.Background;
using Webby.AchievementService.Services.Handlers;
using Webby.NotificationService.GrpcClient;

namespace Webby.AchievementService.Extensions;

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

   public static void AddInterceptors(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddTransient<GrpcExceptionInterceptor>();
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
      serviceCollection.AddScoped<IAchievementRepository, AchievementRepository>();
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IAchievementService, Services.AchievementService>();
      serviceCollection.AddScoped<IStorageService, StorageService>();
      serviceCollection.AddScoped<AchievementHandler>();
   }

   public static void ConfigureGrpcConnection(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddGrpc(options =>
      {
         options.Interceptors.Add<GrpcExceptionInterceptor>();
      });

      serviceCollection.AddGrpcClient<NotificationGrpcService.NotificationGrpcServiceClient>(options =>
      {
         options.Address = new Uri("http://localhost:5007");
      });
      
   }

   public static void AddBackgroundWorkers(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddHostedService<RedisWorker>();
      serviceCollection.AddHostedService<StreamCleanupWorker>();
   }

   public static void AddHelpers(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<INotificationFactory,NotificationFactory>();
   }
   
   public static void ConfigureOptionDependencies(this IServiceCollection serviceCollection, IConfiguration config)
   {
      serviceCollection.Configure<AwsOptions>(config.GetSection(nameof(AwsOptions)));
   }
   
   public static void ConfigureRedisConnection(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddSingleton<IConnectionMultiplexer>(sp =>
         ConnectionMultiplexer.Connect(configuration.GetConnectionString("Redis") ?? string.Empty));
   }
   
}