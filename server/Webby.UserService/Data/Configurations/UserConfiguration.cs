using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class UserConfiguration : IEntityTypeConfiguration<User>
{
   public void Configure(EntityTypeBuilder<User> builder)
   {
      builder.ToTable("Users");

      builder.HasKey(u => u.UserId);
      
      builder.HasMany(u => u.Followers)
         .WithOne()
         .HasForeignKey(uf => uf.UserId)
         .OnDelete(DeleteBehavior.Restrict);
      
      builder.HasMany(u => u.Following)
         .WithOne()
         .HasForeignKey(uf => uf.FollowerId)
         .OnDelete(DeleteBehavior.Restrict);
   }
}