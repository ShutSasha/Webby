using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class UserPremiumConfiguration : IEntityTypeConfiguration<UserPremium>
{
   public void Configure(EntityTypeBuilder<UserPremium> builder)
   {
      builder.ToTable("UserPremiums");

      builder.HasKey(up => up.UserId);

      builder.Property(up => up.ExpiresAt)
         .IsRequired();

      builder.HasOne<User>()
         .WithOne()
         .HasForeignKey<UserPremium>(up => up.UserId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}