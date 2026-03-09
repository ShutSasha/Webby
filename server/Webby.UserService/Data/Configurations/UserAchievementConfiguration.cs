using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class UserAchievementsConfiguration : IEntityTypeConfiguration<UserAchievement>
{
   public void Configure(EntityTypeBuilder<UserAchievement> builder)
   {
      builder.HasKey(ua => new { ua.AchievementId, ua.UserId });

      builder.HasOne<User>()
         .WithMany(u => u.UserAchievements)
         .HasForeignKey(ua => ua.UserId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.HasOne<Achievement>()
         .WithMany(a => a.UserAchievements)
         .HasForeignKey(ua => ua.AchievementId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}