using Microsoft.EntityFrameworkCore;
using Webby.NotificationService.Data.Configurations;
using Webby.NotificationService.Models;

namespace Webby.NotificationService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<Notification> Notifications { get; set; }

   protected override void OnModelCreating(ModelBuilder modelBuilder)
   {
      modelBuilder.ApplyConfiguration(new NotificationConfiguration());
   }
}