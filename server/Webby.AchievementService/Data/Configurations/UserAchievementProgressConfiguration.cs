using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Data.Configurations;

public class UserAchievementProgressConfiguration: IEntityTypeConfiguration<UserAchievementProgress>
{
   public void Configure(EntityTypeBuilder<UserAchievementProgress> builder)
   {
      builder.ToTable("user_achievement_progresses");

      builder.HasKey(uap => new { uap.UserId, uap.AchievementId });

      builder.HasOne(uap => uap.Achievement)
         .WithMany(a => a.UserAchievementsProgresses)
         .HasForeignKey(uap => uap.AchievementId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.Property(uap => uap.CurrentValue)
         .IsRequired()
         .HasDefaultValue(0);
   }
}