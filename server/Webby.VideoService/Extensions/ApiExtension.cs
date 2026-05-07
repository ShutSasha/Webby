using System.Security.Cryptography.X509Certificates;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Data;
using Webby.VideoService.Helpers.External;
using Webby.VideoService.Helpers.Queue;
using Webby.VideoService.Helpers.Storage;
using Webby.VideoService.Interfaces.Helpers;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Repositories;
using Webby.VideoService.Services;
using Webby.VideoService.Services.Background;

namespace Webby.VideoService.Extensions;

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
            Title = "Video Service API",
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

   public static void AddDbConnection(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddDbContext<AppDbContext>(options =>
      {
         options.UseNpgsql(configuration.GetConnectionString(nameof(AppDbContext)));
      });
   }

   public static void AddRepositories(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IPlaylistRepository, PlaylistRepository>();
      serviceCollection.AddScoped<ITagRepository, TagRepository>();
      serviceCollection.AddScoped<IVideoRepository, VideoRepository>();
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<VideoUploadProcessor>();
      serviceCollection.AddSingleton<IBackgroundTaskQueue, BackgroundTaskQueue>();
      serviceCollection.AddHostedService<QueuedHostedService>();
      serviceCollection.AddScoped<IPlaylistService, PlaylistService>();
      serviceCollection.AddScoped<ITagService, TagService>();
      serviceCollection.AddScoped<IVideoService,Services.VideoService>();
      serviceCollection.AddScoped<IStorageService, StorageService>();
      serviceCollection.AddScoped<IStreamService, StreamService>();
      serviceCollection.AddScoped<IYouTubeSearchService,YoutubeSearchService>();
      serviceCollection.AddScoped<ITwitchSearchService, TwitchSearchService>();
   }

   public static void AddExternalServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddHttpClient<YoutubeSearchService>(options =>
      {
         options.BaseAddress = new Uri(DefaultLinks.BaseExternalApiYouTubeUrl);
      });
      
      serviceCollection.AddHttpClient<TwitchSearchService>(options =>
      {
         options.BaseAddress = new Uri(DefaultLinks.BaseExternalApiTwitchLink);
      });
   }

   public static void ConfigureOptionDependencies(this IServiceCollection serviceCollection, IConfiguration config)
   {
      serviceCollection.Configure<AwsOptions>(config.GetSection(nameof(AwsOptions)));
      serviceCollection.Configure<ExternalServicesOptions>(config.GetSection("ExternalServices"));
      
   }

   public static void ConfigureGrpcConnections(this IServiceCollection serviceCollection, IConfiguration configuration)
   {
      serviceCollection.AddGrpcClient<UserGrpcService.UserGrpcServiceClient>(o =>
      {
         o.Address = new Uri(configuration["GrpcClients:UserServiceUrl"]);
      });
   }
}