using Webby.VideoService.Models;

namespace Webby.VideoService.Dtos.Playlist;

public class GetPlaylistResponse
{
   public PlaylistDto Playlist { get; set; }
   public List<Models.Video> Videos { get; set; }
}