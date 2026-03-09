using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class UserAchievementConfiguration : IEntityTypeConfiguration<UserAchievement>
{
   public void Configure(EntityTypeBuilder<UserAchievement> builder)
   {
      builder.HasKey(ua => new { ua.AchievementId, ua.UserId });

      builder.HasOne(ua => ua.User)
         .WithMany(u => u.UserAchievements)
         .HasForeignKey(ua => ua.UserId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.HasOne(ua => ua.Achievement)
         .WithMany(a => a.UserAchievements)
         .HasForeignKey(ua => ua.AchievementId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}