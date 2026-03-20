namespace Webby.VideoService.Models;

public class PlaylistVideo
{
   public Guid PlaylistId { get; set; }
   public Playlist Playlist { get; set; }
   public Guid VideoId { get; set; }
   public Video Video{ get; set; }
}