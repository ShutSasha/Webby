using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class AchievementConfiguration : IEntityTypeConfiguration<Achievement>
{
   public void Configure(EntityTypeBuilder<Achievement> builder)
   {
      builder.ToTable("Achievements");

      builder.HasKey(ac => ac.AchievementId);
      
   }
}