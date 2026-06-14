using Webby.VideoService.Dtos.External;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Helpers;

public interface IExternalContentFetcher
{
   Task<ExternalContentData> FetchExternalContentAsync(IEnumerable<PlaylistVideo> playlistVideos);
   bool IsExternalContentAvailable(PlaylistVideo pv, ExternalContentData data);
}