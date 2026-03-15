namespace Webby.VideoService.Dtos.Playlist;

public class UpdatePlaylistRequest
{
   public Guid PlaylistId { get; set; }
   public string Name { get; set; }
   public string? Description { get; set; }
}