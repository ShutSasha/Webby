using Microsoft.EntityFrameworkCore;
using Webby.VideoService.Data.Configurations;
using Webby.VideoService.Models;

namespace Webby.VideoService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<Video> Videos { get; set; }
   public DbSet<Tag> Tags { get; set; }
   public DbSet<VideoTag> VideoTags { get; set; }
   public DbSet<Playlist> Playlists { get; set; }
   public DbSet<PlaylistVideo> PlaylistVideos { get; set; }
   public DbSet<UserView> UserViews { get; set; }
   
   protected override void OnModelCreating(ModelBuilder modelBuilder)
   {
      modelBuilder.ApplyConfiguration(new VideoConfiguration());
      modelBuilder.ApplyConfiguration(new TagConfiguration());
      modelBuilder.ApplyConfiguration(new VideoTagConfiguration());
      modelBuilder.ApplyConfiguration(new PlaylistConfiguration());
      modelBuilder.ApplyConfiguration(new PlaylistVideoConfiguration());
      modelBuilder.ApplyConfiguration(new UserViewConfiguration());
   }
}