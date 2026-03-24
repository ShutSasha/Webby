namespace Webby.VideoService.Dtos.Playlist;

public class PlaylistDto
{
   public Guid PlaylistId { get; set; }
   public Guid UserId { get; set; }
   public string Name { get; set; }
   public int CountOfVideos { get; set; }
   public bool IsPrivate { get; set; }
   public string PlaylistCover { get; set; }
}