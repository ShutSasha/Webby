using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Models;

public class PlaylistVideo
{
   public Guid PlaylistVideoId { get; set; }

   public Guid PlaylistId { get; set; }
   public Playlist Playlist { get; set; } = null!;
   
   public SystemPlatforms Platform { get; set; }
   
   public MediaType MediaType { get; set; }
   public Guid? InternalContentId { get; set; }
   public Video? Video { get; set; }
   
   public string? ExternalContentId { get; set; }
   public DateTime CreatedAt { get; set; }
}