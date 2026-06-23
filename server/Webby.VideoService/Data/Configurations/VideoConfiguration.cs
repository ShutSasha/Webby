using Webby.VideoService.Models;

namespace Webby.VideoService.Data.Configurations;

using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

public class VideoConfiguration : IEntityTypeConfiguration<Video>
{
   public void Configure(EntityTypeBuilder<Video> builder)
   {
      builder.ToTable("Videos");

      builder.HasKey(v => v.VideoId);

      builder.Property(v => v.Name)
         .IsRequired()
         .HasMaxLength(500);

      builder.Property(v => v.Duration)
         .HasDefaultValue(0L);
      
      builder.Property(v => v.Description)
         .HasMaxLength(1000);

      builder.Property(v => v.Views)
         .HasDefaultValue(0);

      builder.Property(v => v.VideoUploadStatus)
         .HasConversion<string>()
         .IsRequired();
      
      builder.Property(v => v.IsPrivate)
         .HasDefaultValue(false);

      builder.Property(v => v.IsBanned)
         .HasDefaultValue(false);
      
      builder.HasMany(v => v.VideoTags)
         .WithOne(vt => vt.Video)
         .HasForeignKey(vt => vt.VideoId);
   }
}
