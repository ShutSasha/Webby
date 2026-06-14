namespace Webby.VideoService.Dtos.Playlist;

public class AddStreamToPlaylistRequest
{
   public Guid PlaylistId { get; set; }
   public List<string> StreamIds { get; set; }
}