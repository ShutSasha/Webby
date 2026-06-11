using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Stream.Enums;
using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Helpers.Converters;
public record struct PlatformParseResult<TEnum>(TEnum Platform, string ActualId) 
    where TEnum : struct, Enum;

public static class PlatformPrefixToPlatformConverter
{

    private static readonly IReadOnlyDictionary<string, VideoPlatform> VideoPlatformMap =
        new Dictionary<string, VideoPlatform>
        {
            { PlatformPrefixesConstants.WebbyPrefix, VideoPlatform.Webby },
            { PlatformPrefixesConstants.YouTubePrefix, VideoPlatform.YouTube }
        };

    private static readonly IReadOnlyDictionary<string, SearchVideoPlatforms> SearchVideoPlatformMap =
        new Dictionary<string, SearchVideoPlatforms>
        {
            { PlatformPrefixesConstants.WebbyPrefix, SearchVideoPlatforms.Webby },
            { PlatformPrefixesConstants.YouTubePrefix, SearchVideoPlatforms.YouTube }
        };

    private static readonly IReadOnlyDictionary<string, SearchStreamPlatforms> StreamPlatformMap =
        new Dictionary<string, SearchStreamPlatforms>
        {
            { PlatformPrefixesConstants.TwitchPrefix, SearchStreamPlatforms.Twitch }
        };

    private static readonly IReadOnlyDictionary<string, SystemPlatforms> SystemPlatformsMap =
        new Dictionary<string, SystemPlatforms>
        {
            { PlatformPrefixesConstants.WebbyPrefix, SystemPlatforms.Webby },
            { PlatformPrefixesConstants.YouTubePrefix, SystemPlatforms.YouTube },
            { PlatformPrefixesConstants.TwitchPrefix, SystemPlatforms.Twitch }
        };
    

    private static PlatformParseResult<TEnum>? Parse<TEnum>(
        string prefixedId,
        IReadOnlyDictionary<string, TEnum> prefixMap)
        where TEnum : struct, Enum
    {
        if (string.IsNullOrWhiteSpace(prefixedId))
            return null;

        foreach (var (prefix, platform) in prefixMap)
        {
            if (prefixedId.StartsWith(prefix, StringComparison.OrdinalIgnoreCase))
            {
                var actualId = prefixedId[prefix.Length..];
                return new PlatformParseResult<TEnum>(platform, actualId);
            }
        }

        return null;
    }

    public static PlatformParseResult<VideoPlatform>? ParseVideoPlatform(string prefixedId)
        => Parse(prefixedId, VideoPlatformMap);

    public static PlatformParseResult<SearchVideoPlatforms>? ParseSearchVideoPlatform(string prefixedId)
        => Parse(prefixedId, SearchVideoPlatformMap);

    public static PlatformParseResult<SearchStreamPlatforms>? ParseStreamPlatform(string prefixedId)
        => Parse(prefixedId, StreamPlatformMap);

    public static PlatformParseResult<SystemPlatforms>? ParseSystemPlatform(string prefixedId)
        => Parse(prefixedId, SystemPlatformsMap);
    
}