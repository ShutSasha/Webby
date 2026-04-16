namespace Webby.VideoService.Constants;

public static class DefaultLinks
{
   public const string PlaylistEmptyLink = "https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/playlists/empty-playlist.png";
   
   public const string BaseExternalApiYouTubeUrl = "https://www.googleapis.com/";
   public const string BaseYouTubeSearchLink = "https://www.googleapis.com/youtube/v3/search";
   public const string BaseYouTubeVideosLink = "https://www.googleapis.com/youtube/v3/videos";
   public const string BaseYouTubeUserIcon = "https://upload.wikimedia.org/wikipedia/commons/thumb/0/09/YouTube_full-color_icon_%282017%29.svg/960px-YouTube_full-color_icon_%282017%29.svg.png";
   public const string BaseYouTubeUserChannelsLinks = "https://www.googleapis.com/youtube/v3/channels";
   
   public const string BaseExternalApiTwitchLink = "https://api.twitch.tv/helix/";
   public const string BaseTwitchStreamLink = "https://api.twitch.tv/helix/streams";
   public const string BaseTwitchUserLink = "https://api.twitch.tv/helix/users";
   public const string BaseTwitchOAuthLink = "https://id.twitch.tv/oauth2/token";
   
   public const string DefaultTwitchPlayerWatchLink = "https://www.twitch.tv/";
}