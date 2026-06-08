namespace Webby.VideoService.Dtos.Playlist;

public class AddStreamToPlaylistsRequest
{
   public List<Guid> PlaylistIds { get; set; }
   public List<string> StreamIds { get; set; }
}