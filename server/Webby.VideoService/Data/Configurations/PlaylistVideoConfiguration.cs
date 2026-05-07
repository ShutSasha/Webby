using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.VideoService.Models;

namespace Webby.VideoService.Data.Configurations;

public class PlaylistVideoConfiguration : IEntityTypeConfiguration<PlaylistVideo>
{
   public void Configure(EntityTypeBuilder<PlaylistVideo> builder)
   {
      builder.ToTable("PlaylistVideos");
      
      builder.HasKey(pv => pv.PlaylistVideoId);
      
      builder.HasOne(pv => pv.Playlist)
         .WithMany(p => p.PlaylistVideos)
         .HasForeignKey(pv => pv.PlaylistId)
         .OnDelete(DeleteBehavior.Cascade);
      
      builder.HasOne(pv => pv.Video)
         .WithMany()
         .HasForeignKey(pv => pv.VideoId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.Property(pv => pv.VideoPlatform)
         .HasConversion<string>()
         .IsRequired();
      
      builder.HasIndex(pv => pv.VideoId);
      builder.HasIndex(pv => pv.ExternalVideoId);
      
      builder.HasIndex(pv => new { pv.PlaylistId, pv.VideoId, pv.ExternalVideoId }).IsUnique();
   }
}