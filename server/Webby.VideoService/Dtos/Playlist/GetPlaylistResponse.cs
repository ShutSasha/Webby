using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Models;

namespace Webby.VideoService.Dtos.Playlist;

public class GetPlaylistResponse
{
   public PlaylistDto Playlist { get; set; }
   public List<VideoDto> Videos { get; set; }
}