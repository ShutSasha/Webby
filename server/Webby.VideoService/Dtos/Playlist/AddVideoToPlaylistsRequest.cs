namespace Webby.VideoService.Dtos.Playlist;

public class AddVideoToPlaylistsRequest
{
   public List<Guid> PlaylistIds { get; set; }
   public List<string> VideoIds { get; set; }
}