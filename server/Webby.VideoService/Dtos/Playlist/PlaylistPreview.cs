namespace Webby.VideoService.Dtos.Playlist;

public class PlaylistPreviewDto
{
   public Guid PlaylistId { get; set; }
   public required string Name { get; set; }
   public required string PlaylistCover { get; set; }
   public required int CountOfVideos { get; set; }
   public required bool IsVideoAdded { get; set; } = false;
}