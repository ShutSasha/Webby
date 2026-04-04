using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.VideoService.Models;

namespace Webby.VideoService.Data.Configurations;

public class UserViewConfiguration : IEntityTypeConfiguration<UserView>
{
   public void Configure(EntityTypeBuilder<UserView> builder)
   {
      builder.ToTable("UserViews");

      builder.HasKey(uv => new { uv.VideoId, uv.UserId });
      
   }
}