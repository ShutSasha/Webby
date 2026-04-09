using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data.Configurations;
using Webby.UserService.Models;

namespace Webby.UserService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<User> Users { get; set; }
   public DbSet<Complaint> Complaints { get; set; }
   public DbSet<UserFollower> UserFollowers { get; set; }
   public DbSet<UserPremium?> UserPremiums { get; set; }
   public DbSet<Achievement> Achievements { get; set; }
   public DbSet<UserAchievement> UserAchievements { get; set; }
   public DbSet<Payment> Payments { get; set; }

   protected override void OnModelCreating(ModelBuilder modelBuilder)
   {
      modelBuilder.ApplyConfiguration(new UserConfiguration());
      modelBuilder.ApplyConfiguration(new ComplaintConfiguration());
      modelBuilder.ApplyConfiguration(new UserFollowerConfiguration());
      modelBuilder.ApplyConfiguration(new UserPremiumConfiguration());
      modelBuilder.ApplyConfiguration(new AchievementConfiguration());
      modelBuilder.ApplyConfiguration(new UserAchievementConfiguration());
      modelBuilder.ApplyConfiguration(new PaymentConfiguration());
   }
}