namespace Webby.VideoService.Models;

public class Playlist
{
   public Guid PlaylistId { get; set; }
   public string Name { get; set; }
   public string? Description { get; set; }
   public Guid UserId { get; set; }
   public ICollection<PlaylistVideo> PlaylistVideos { get; set; }
}