using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.VideoService.Models;

namespace Webby.VideoService.Data.Configurations;

public class PlaylistConfiguration : IEntityTypeConfiguration<Playlist>
{
   public void Configure(EntityTypeBuilder<Playlist> builder)
   {
      builder.ToTable("Playlists");

      builder.HasKey(p => p.PlaylistId);

      builder.Property(p => p.Name)
         .IsRequired()
         .HasMaxLength(500);

      builder.Property(p => p.Description)
         .HasMaxLength(1000);

      builder.Property(p => p.UserId)
         .IsRequired();

      builder.HasIndex(p => p.UserId);

      builder.HasMany(p => p.PlaylistVideos)
         .WithOne(pv => pv.Playlist)
         .HasForeignKey(pv => pv.PlaylistId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}