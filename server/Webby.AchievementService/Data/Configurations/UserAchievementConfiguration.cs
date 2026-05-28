using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Data.Configurations;

public class UserAchievementConfiguration : IEntityTypeConfiguration<UserAchievement>
{
   public void Configure(EntityTypeBuilder<UserAchievement> builder)
   {
      //TODO: rename to 'user_achievements' after implementing achievement system logic
      builder.ToTable("UserAchievements");
      
      builder.HasKey(ua => new { ua.AchievementId, ua.UserId });

      builder.HasOne(ua => ua.Achievement)
         .WithMany(a => a.UserAchievements)
         .HasForeignKey(ua => ua.AchievementId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}