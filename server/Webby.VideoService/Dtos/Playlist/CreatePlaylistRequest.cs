using System.ComponentModel.DataAnnotations;

namespace Webby.VideoService.Dtos.Playlist;

public class CreatePlaylistRequest
{
   [StringLength(200, MinimumLength = 3, ErrorMessage = "Field {0} must be less than {1} characters and greater than {2} characters")]
   public required string Name { get; set; }
   public required bool IsPrivate { get; set; } = false;
}