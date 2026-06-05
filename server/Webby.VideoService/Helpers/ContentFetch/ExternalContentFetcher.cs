using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Interfaces.Helpers;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;


public class ExternalContentFetcher : IExternalContentFetcher
{
    private readonly IYouTubeSearchService _youtubeSearchService;
    private readonly ITwitchSearchService _twitchSearchService;

    public ExternalContentFetcher(IYouTubeSearchService youtubeSearchService, ITwitchSearchService twitchSearchService)
    {
        _youtubeSearchService = youtubeSearchService;
        _twitchSearchService = twitchSearchService;
    }

    public async Task<ExternalContentData> FetchExternalContentAsync(IEnumerable<PlaylistVideo> playlistVideos)
    {
        var result = new ExternalContentData();

        var youtubeIds = playlistVideos
            .Where(pv => pv is { MediaType: MediaType.Video, Platform: SystemPlatforms.YouTube } && !string.IsNullOrEmpty(pv.ExternalContentId))
            .Select(pv => pv.ExternalContentId!)
            .ToList();

        var twitchIds = playlistVideos
            .Where(pv => pv is { MediaType: MediaType.LiveStream, Platform: SystemPlatforms.Twitch } && !string.IsNullOrEmpty(pv.ExternalContentId))
            .Select(pv => pv.ExternalContentId!)
            .ToList();

        if (youtubeIds.Count > 0)
        {
            var ytList = await _youtubeSearchService.GetList(youtubeIds);
            result.YoutubeVideos = ytList.ToDictionary(v => v.VideoId!);
            result.UnavailableCount += youtubeIds.Count - result.YoutubeVideos.Count;
        }

        if (twitchIds.Count > 0)
        {
            var twitchList = await _twitchSearchService.GetList(twitchIds);
            result.TwitchStreams = twitchList.ToDictionary(s => s.StreamerId!);
            result.UnavailableCount += twitchIds.Count - result.TwitchStreams.Count;
        }

        return result;
    }

    public bool IsExternalContentAvailable(PlaylistVideo pv, ExternalContentData data)
    {
        if (string.IsNullOrEmpty(pv.ExternalContentId)) return false;

        if (pv.MediaType == MediaType.Video && pv.Platform == SystemPlatforms.YouTube)
        {
            return data.YoutubeVideos.ContainsKey(PlatformPrefixesConstants.YouTubePrefix + pv.ExternalContentId);
        }

        if (pv.MediaType == MediaType.LiveStream && pv.Platform == SystemPlatforms.Twitch)
        {
            return data.TwitchStreams.ContainsKey(PlatformPrefixesConstants.TwitchPrefix + pv.ExternalContentId);
        }

        return false;
    }
}