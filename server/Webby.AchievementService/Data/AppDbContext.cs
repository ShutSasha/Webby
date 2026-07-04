using Microsoft.EntityFrameworkCore;
using Webby.AchievementService.Data.Configurations;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<Achievement> Achievements { get; set; }
   public DbSet<UserAchievement> UserAchievements { get; set; }
   public DbSet<UserAchievementProgress> UserAchievementProgresses { get; set; }
   public DbSet<EventType> EventTypes { get; set; }

   protected override void OnModelCreating(ModelBuilder modelBuilder)
   {
      modelBuilder.ApplyConfiguration(new AchievementConfiguration());
      modelBuilder.ApplyConfiguration(new UserAchievementConfiguration());
      modelBuilder.ApplyConfiguration(new UserAchievementProgressConfiguration());
   }
}