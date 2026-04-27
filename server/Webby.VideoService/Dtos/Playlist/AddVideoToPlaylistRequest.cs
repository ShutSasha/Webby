using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Dtos.Playlist;

public class AddVideoToPlaylistRequest
{
   public Guid PlaylistId { get; set; }
   public List<AddVideoToPlaylistItem> VideoItems { get; set; }
}

public class AddVideoToPlaylistItem
{
   public VideoPlatform VideoPlatform { get; set; }
   public string? ItemId { get; set; }
}