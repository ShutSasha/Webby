using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Dtos.Playlist;

public class AddVideoToPlaylistRequest
{
   public Guid PlaylistId { get; set; }
   public List<string> VideoIds { get; set; }
}

public record PlaylistItemInput(SystemPlatforms VideoPlatform, string ItemId);