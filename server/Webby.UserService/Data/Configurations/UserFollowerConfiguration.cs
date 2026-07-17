using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class UserFollowerConfiguration : IEntityTypeConfiguration<UserFollower>
{
   public void Configure(EntityTypeBuilder<UserFollower> builder)
   {
      builder.ToTable("UserFollowers");

      builder.HasKey(uf => new { uf.UserId, uf.FollowerId });

      builder.HasIndex(uf => uf.FollowerId);

      builder.HasOne(uf => uf.FollowedUser)
         .WithMany(u => u.Followers)
         .HasForeignKey(uf => uf.UserId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.HasOne(uf => uf.FollowerUser)
         .WithMany(u => u.Following)
         .HasForeignKey(uf => uf.FollowerId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}