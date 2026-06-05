using System.Net.Http.Headers;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Options;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.External;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services;

public class TwitchSearchService : ITwitchSearchService
{
    private readonly HttpClient _httpClient;
    private readonly TwitchOptions _options;
    private string? _accessToken;

    public TwitchSearchService(HttpClient httpClient, IOptions<ExternalServicesOptions> options)
    {
        _httpClient = httpClient;
        _options = options.Value.TwitchOptions;
    }
    
    public async Task<PagedResponse<StreamDto>> SearchAsync(string? searchText, int pageSize, int page, string? nextPageToken)
    {
        await AuthenticateAsync();

        bool isSearch = !string.IsNullOrWhiteSpace(searchText);
        string? NewNextPageToken = null;
        List<TwitchItem> streamData = new();

        if (isSearch)
        {
            var searchUrl = DefaultLinks.BaseTwitchUserLink;
            var searchParams = new Dictionary<string, string?>
            {
                ["query"] = searchText,
                ["live_only"] = "true",
                ["first"] = pageSize.ToString()
            };

            if (!string.IsNullOrEmpty(nextPageToken))
                searchParams["after"] = nextPageToken;

            var searchResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(
                QueryHelpers.AddQueryString(searchUrl, searchParams));

            if (searchResponse == null)
            {
                return new PagedResponse<StreamDto>();
            }

            NewNextPageToken = searchResponse.Pagination?.Cursor;

            if (searchResponse.Data.Any())
            {
                var ids = searchResponse.Data.Select(i => i.Id);
                var streamsUrl = DefaultLinks.BaseTwitchStreamLink + $"?user_id={string.Join("&user_id=", ids)}";
                
                var streamsResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(streamsUrl);
                streamData = streamsResponse.Data;
            }
        }
        else
        {
            var topUrl = DefaultLinks.BaseTwitchStreamLink;
            var topParams = new Dictionary<string, string?>
            {
                ["first"] = pageSize.ToString()
            };

            if (!string.IsNullOrEmpty(nextPageToken))
                topParams["after"] = nextPageToken;

            var topResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(
                QueryHelpers.AddQueryString(topUrl, topParams));

            streamData = topResponse.Data;
            NewNextPageToken = topResponse.Pagination?.Cursor;
        }

        if (!streamData.Any()) 
            return new PagedResponse<StreamDto> { Items = [] };
        
        var userIds = streamData.Select(i => i.UserId ?? i.Id).Distinct();
        var avatars = await GetUsersAvatarsAsync(userIds);
        
        var resultItems = streamData.Select(item => new StreamDto
        {
            StreamerId = PlatformPrefixesConstants.TwitchPrefix + item.UserId,
            Name = item.Title,
            StartedAt = item.StartedAt,
            PreviewUrl = item.ThumbnailUrl.Replace("{width}",_options.ThumbnailWidth).Replace("{height}", _options.ThumbnailHeight),
            StreamUrl = DefaultLinks.DefaultTwitchPlayerWatchLink + item.UserName ?? item.BroadcasterLogin!,
            Viewers = item.ViewerCount,
            Source = "Twitch",
            StreamerInformation =new UserStreamDto
            {
                Username = item.UserName ?? item.DisplayName ?? "Unknown",
                AvatarUrl = avatars.GetValueOrDefault(item.UserId ?? item.Id) ?? "",
            }
        }).ToList();

        return new PagedResponse<StreamDto>
        {
            Items = resultItems,
            NextPageToken = NewNextPageToken,
            PageSize = pageSize,
            Page = page
        };
    }

    public async Task<StreamDto> FindById(string streamId)
    {
        if (string.IsNullOrEmpty(streamId)) 
            return null;

        await AuthenticateAsync();
        
        var streamUrl = $"{DefaultLinks.BaseTwitchStreamLink}?user_id={streamId}";
        var streamResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(streamUrl);
        var streamItem = streamResponse.Data?.FirstOrDefault();

        if (streamItem != null)
        {
            var avatars = await GetUsersAvatarsAsync([streamItem.UserId]);
            return MapToStreamDto(streamItem, avatars.GetValueOrDefault(streamItem.UserId));
        }
        
        var videoUrl = $"https://api.twitch.tv/helix/videos?id={streamId}";
        var vResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(videoUrl);
        var videoItem = vResponse != null ? vResponse.Data?.FirstOrDefault() : null;

        if (videoItem != null)
        {
            var avatars = await GetUsersAvatarsAsync([videoItem.UserId]);
            return MapToStreamDto(videoItem, avatars.GetValueOrDefault(videoItem.UserId));
        }

        throw new ApiException("Get twitch stream error", 404, "Stream wasn't found");
    }

    public async Task<List<StreamDto>> GetList(List<string> sourceIds)
    {
        if (sourceIds == null || sourceIds.Count == 0)
            return new List<StreamDto>();

        await AuthenticateAsync();

        var results = new List<StreamDto>();
        var items = new List<TwitchItem>();

        foreach (var chunk in sourceIds.Chunk(100))
        {
            var streamUrl = $"{DefaultLinks.BaseTwitchStreamLink}?user_id={string.Join("&user_id=", chunk)}";
            var streamResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(streamUrl);
            
            var foundUserIds = new HashSet<string>();

            if (streamResponse?.Data != null)
            {
                items.AddRange(streamResponse.Data);
                foreach (var stream in streamResponse.Data)
                {
                    if (!string.IsNullOrEmpty(stream.UserId))
                    {
                        foundUserIds.Add(stream.UserId);
                    }
                }
            }

            var missingIds = chunk.Except(foundUserIds).ToList();

            if (missingIds.Count > 0)
            {
                var videoUrl = $"https://api.twitch.tv/helix/videos?id={string.Join("&id=", missingIds)}";
                var videoResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(videoUrl);
                
                if (videoResponse?.Data != null)
                {
                    items.AddRange(videoResponse.Data);
                }
            }
        }

        if (items.Count == 0)
            return results;

        var userIds = items.Select(i => i.UserId ?? i.Id).Where(id => !string.IsNullOrEmpty(id)).Distinct().ToList();
        var avatars = new Dictionary<string, string>();

        foreach (var userChunk in userIds.Chunk(100))
        {
            var chunkAvatars = await GetUsersAvatarsAsync(userChunk);
            foreach (var kvp in chunkAvatars)
            {
                avatars[kvp.Key] = kvp.Value;
            }
        }

        foreach (var item in items)
        {
            results.Add(MapToStreamDto(item, avatars.GetValueOrDefault(item.UserId ?? item.Id!)));
        }

        return results;
    }

    private StreamDto MapToStreamDto(TwitchItem? item, string? avatarUrl)
    {
        return new StreamDto
        {
            StreamerId= PlatformPrefixesConstants.TwitchPrefix + item.UserId,
            Name = item.Title,
            StartedAt = item.StartedAt,
            PreviewUrl = item.ThumbnailUrl.Replace("{width}", _options.ThumbnailWidth).Replace("{height}", _options.ThumbnailHeight),
            StreamUrl = DefaultLinks.DefaultTwitchPlayerWatchLink + item.UserName ?? item.BroadcasterLogin!,
            Viewers = item.ViewerCount,
            Source = "Twitch",
            StreamerInformation= new UserStreamDto
            {
                Username = item.UserName ?? item.DisplayName ?? "Unknown",
                AvatarUrl = avatarUrl ?? "",
            }
        };
    }

    private async Task<Dictionary<string, string>> GetUsersAvatarsAsync(IEnumerable<string> userIds)
    {
        var query = string.Join("&id=", userIds);
        var url = DefaultLinks.BaseTwitchUserLink + $"?id={query}";
        
        var response = await SendTwitchRequest<TwitchResponse<TwitchUser>>(url);
        return response.Data.ToDictionary(u => u.Id, u => u.ProfileImageUrl);
    }
    
    private async Task<T> SendTwitchRequest<T>(string url)
    {
        using var request = new HttpRequestMessage(HttpMethod.Get, url);
        request.Headers.Add("Client-Id", _options.ClientId);
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", _accessToken);
        
        var response = await _httpClient.SendAsync(request);

        if (!response.IsSuccessStatusCode)
        {
            return default!;
        }
        
        return await response.Content.ReadFromJsonAsync<T>();
    }
    
    private async Task AuthenticateAsync()
    {
        if (!string.IsNullOrEmpty(_accessToken)) return;

        var authUrl = DefaultLinks.BaseTwitchOAuthLink +
                      $"?client_id={_options.ClientId}" +
                      $"&client_secret={_options.ClientSecret}" +
                      "&grant_type=client_credentials";

        var response = await _httpClient.PostAsync(authUrl, null);
        response.EnsureSuccessStatusCode();

        var tokenData = await response.Content.ReadFromJsonAsync<TwitchTokenResponse>();
        _accessToken = tokenData?.AccessToken;
    }
}