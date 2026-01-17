using Microsoft.EntityFrameworkCore;
using Webby.AuthService.Data;
using Webby.AuthService.Helpers;
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
   }

   public static void AddServices(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IMailService, MailService>();
      serviceCollection.AddScoped<IAuthService, Services.AuthService>();
   }

   public static void AddHelpers(this IServiceCollection serviceCollection)
   {
      serviceCollection.AddScoped<IPasswordHasher, PasswordHasher>();
   }
   
}