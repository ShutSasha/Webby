using System.Text;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using Webby.ApiGetaway.Helpers.Exception;
using Webby.ApiGetaway.Helpers.Jwt;
using Ocelot.Configuration.File;
using Ocelot.DependencyInjection;

namespace Webby.ApiGetaway.Extensions;

public static class ApiExtension
{
   public static void AddCorsPolicy(this IServiceCollection serviceCollection, string policyName)
   {
      serviceCollection.AddCors(corsOptions =>
      {
         corsOptions.AddPolicy(policyName, policy =>
         {
            policy
               .WithOrigins("http://localhost:3000")
               .AllowAnyMethod()
               .AllowCredentials()
               .AllowAnyHeader();
         });
      });
   }

   public static void AddSwaggerInfo(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddSwaggerGen(c =>
      {
         c.SwaggerDoc("v1", new OpenApiInfo { Title = "API Gateway", Version = "v1" });
      });

   }
   
   public static void AddJwtAuthentication(this IServiceCollection services, IConfiguration configuration)
   {
      var jwtOptions = configuration
         .GetSection("JwtOptions")
         .Get<JwtOptions>();

      services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
         .AddJwtBearer(options =>
         {
            options.TokenValidationParameters = new TokenValidationParameters
            {
               ValidateIssuer = false,
               ValidateAudience = false,
               ValidateLifetime = true,
               ValidateIssuerSigningKey = true,

               IssuerSigningKey = new SymmetricSecurityKey(
                  Encoding.UTF8.GetBytes(jwtOptions!.AccessSecretKey)),
               ClockSkew = TimeSpan.Zero
            };
            
            options.Events = new JwtBearerEvents
            {
               OnChallenge = async context =>
               {
                  context.HandleResponse();

                  context.Response.StatusCode = StatusCodes.Status401Unauthorized;
                  context.Response.ContentType = "application/json";

                  var response = new ApiException("You're not logged in. Sign in your account or sign up", 401);

                  await context.Response.WriteAsJsonAsync(response);
               }
            };
         });

      services.AddAuthorization();
   }
   
   public static void RegisterApiConfig(this ConfigurationManager configuration, IWebHostEnvironment env)
   {
      configuration.SetBasePath(env.ContentRootPath)
         .AddOcelot($"Configurations/{env.EnvironmentName}", env);
   }
   
   public static IApplicationBuilder UseApiExceptionHandling(this IApplicationBuilder app)
   {
      return app.Use(async (context, next) =>
      {
         await next();

         if (context.Response.HasStarted)
            return;

         if (context.Response.StatusCode == StatusCodes.Status401Unauthorized)
         {
            context.Response.ContentType = "application/json";

            var response = new ApiException("You're not logged in. Sign in your account or sign up",401);

            await context.Response.WriteAsJsonAsync(response);
         }

         if (context.Response.StatusCode == StatusCodes.Status403Forbidden)
         {
            context.Response.ContentType = "application/json";

            var response = new ApiException("You don't have permission",403);

            await context.Response.WriteAsJsonAsync(response);
         }
      });
   }
   
}