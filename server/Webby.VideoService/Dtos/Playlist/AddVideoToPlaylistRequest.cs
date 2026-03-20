namespace Webby.VideoService.Dtos.Playlist;

public class AddVideoToPlaylistRequest
{
   public Guid PlaylistId { get; set; }
   public List<Guid> VideoIds { get; set; }
}