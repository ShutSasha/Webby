using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Models;

namespace Webby.VideoService.Dtos.Playlist;

public class GetPlaylistResponse
{
   public PlaylistDto Playlist { get; set; }
   public VideoDto? FirstVideo { get; set; }
   public int HiddenVideosCount { get; set; }
}