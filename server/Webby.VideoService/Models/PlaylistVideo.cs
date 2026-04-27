using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Models;

public class PlaylistVideo
{
   public Guid PlaylistVideoId { get; set; }
   public Guid PlaylistId { get; set; }
   public Playlist Playlist { get; set; }
   public VideoPlatform VideoPlatform { get; set; }
   public Guid? VideoId { get; set; }
   public Video? Video { get; set; }
   public string? ExternalVideoId { get; set; }

   public DateTime CreatedAt { get; set; }
}