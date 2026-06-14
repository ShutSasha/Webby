using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.VideoService.Models;

namespace Webby.VideoService.Data.Configurations;

public class PlaylistVideoConfiguration : IEntityTypeConfiguration<PlaylistVideo>
{
   public void Configure(EntityTypeBuilder<PlaylistVideo> builder)
   {
      builder.ToTable("PlaylistVideos");

      builder.HasKey(x => x.PlaylistVideoId);

      builder.HasOne(x => x.Playlist)
         .WithMany(x => x.PlaylistVideos)
         .HasForeignKey(x => x.PlaylistId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.Property(x => x.Platform)
         .HasConversion<string>()
         .IsRequired();

      builder.Property(x => x.MediaType)
         .HasConversion<string>()
         .IsRequired();
      
      builder.HasOne(x => x.Video)
         .WithMany()
         .HasForeignKey(x => x.InternalContentId)
         .OnDelete(DeleteBehavior.Cascade);

      builder.Property(x => x.ExternalContentId)
         .HasMaxLength(255);

      builder.HasIndex(x => x.PlaylistId);

      builder.HasIndex(x => x.InternalContentId);

      builder.HasIndex(x => x.ExternalContentId);

      builder.HasIndex(x => new
      {
         x.PlaylistId,
         x.InternalContentId,
         x.ExternalContentId
      }).IsUnique();
   }
}