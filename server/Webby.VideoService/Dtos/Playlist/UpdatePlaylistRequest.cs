namespace Webby.VideoService.Dtos.Playlist;

public class UpdatePlaylistRequest : CreatePlaylistRequest
{
   public Guid PlaylistId { get; set; }
}