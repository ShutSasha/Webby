namespace Webby.VideoService.Dtos.Playlist;

public class CreatePlaylistRequest
{
   public required string Name { get; set; }
   public required string? Description { get; set; }
   public required bool IsPrivate { get; set; } = false;
}