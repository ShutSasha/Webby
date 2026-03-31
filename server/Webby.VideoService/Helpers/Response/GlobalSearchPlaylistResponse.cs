using Webby.VideoService.Dtos.Playlist;

namespace Webby.VideoService.Helpers.Response;

public class GlobalSearchPlaylistResponse
{
   public SearchSection<SearchPlaylistDto> WebbyPlaylists { get; set; } = new();
}