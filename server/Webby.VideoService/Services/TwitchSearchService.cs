using System.Net.Http.Headers;
using System.Text.Json;
using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Extensions.Options;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.External;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services;

public class TwitchSearchService : IExternalVideoSearchService
{
    private readonly HttpClient _httpClient;
    private readonly TwitchOptions _options;
    private string? _accessToken;

    public TwitchSearchService(HttpClient httpClient, IOptions<ExternalServicesOptions> options)
    {
        _httpClient = httpClient;
        _options = options.Value.TwitchOptions;
    }
    
    public async Task<PagedResponse<VideoDto>> SearchAsync(SearchVideoOptions options)
{
    await AuthenticateAsync();

    bool isSearch = !string.IsNullOrWhiteSpace(options.SearchText);
    string? nextPageToken = null;
    List<TwitchItem> streamData = new();

    if (isSearch)
    {
        var searchUrl = DefaultLinks.BaseTwitchUserLink;
        var searchParams = new Dictionary<string, string?>
        {
            ["query"] = options.SearchText,
            ["live_only"] = "true",
            ["first"] = options.PageSize.ToString()
        };
        
        if (!string.IsNullOrEmpty(options.NextPageToken))
            searchParams["after"] = options.NextPageToken;

        var searchResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(
            QueryHelpers.AddQueryString(searchUrl, searchParams));

        if (searchResponse == null)
        {
            return new PagedResponse<VideoDto>();
        }

        nextPageToken = searchResponse.Pagination?.Cursor;

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
            ["first"] = options.PageSize.ToString()
        };

        if (!string.IsNullOrEmpty(options.NextPageToken))
            topParams["after"] = options.NextPageToken;

        var topResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(
            QueryHelpers.AddQueryString(topUrl, topParams));

        streamData = topResponse.Data;
        nextPageToken = topResponse.Pagination?.Cursor;
    }

    if (!streamData.Any()) 
        return new PagedResponse<VideoDto> { Items = new List<VideoDto>() };
    
    var userIds = streamData.Select(i => i.UserId ?? i.Id).Distinct();
    var avatars = await GetUsersAvatarsAsync(userIds);
    
    var resultItems = streamData.Select(item => new VideoDto
    {
        VideoId = item.Id,
        Name = item.Title,
        Description = $"Streaming: {item.GameName}",
        Source = "Twitch",
        PreviewUrl = item.ThumbnailUrl.Replace("{width}", "640").Replace("{height}", "360"),
        VideoUrl = $"https://www.twitch.tv/{item.UserName ?? item.BroadcasterLogin}",
        CreatedAt = item.StartedAt,
        IsPrivate = false,
        VideoUploadStatus = VideoStatus.Ready,
        Duration = 0L,
        Views = item.ViewerCount,
        User = new UserVideoDto
        {
            UserId = item.UserId,
            Username = item.UserName ?? item.DisplayName ?? "Unknown",
            AvatarUrl = avatars.GetValueOrDefault(item.UserId ?? item.Id) ?? "",
            IsFollowed = false
        }
    }).ToList();

    return new PagedResponse<VideoDto>
    {
        Items = resultItems,
        NextPageToken = nextPageToken,
        PageSize = options.PageSize
    };
}

    public async Task<VideoDto> FindById(string videoId)
    {
        if (string.IsNullOrEmpty(videoId)) return null;

        await AuthenticateAsync();
        
        var streamUrl = $"{DefaultLinks.BaseTwitchStreamLink}?user_id={videoId}";
        var streamResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(streamUrl);
        var streamItem = streamResponse.Data?.FirstOrDefault();

        if (streamItem != null)
        {
            var avatars = await GetUsersAvatarsAsync(new[] { streamItem.UserId });
            return MapToVideoDto(streamItem, avatars.GetValueOrDefault(streamItem.UserId));
        }
        
        var videoUrl = $"https://api.twitch.tv/helix/videos?id={videoId}";
        var vResponse = await SendTwitchRequest<TwitchResponse<TwitchItem>>(videoUrl);
        var videoItem = vResponse.Data?.FirstOrDefault();

        if (videoItem != null)
        {
            var avatars = await GetUsersAvatarsAsync(new[] { videoItem.UserId });
            return MapToVideoDto(videoItem, avatars.GetValueOrDefault(videoItem.UserId));
        }

        throw new ApiException("Get twitch stream error", 404, "Stream wasn't found");
    }
    
    private VideoDto MapToVideoDto(TwitchItem item, string? avatarUrl)
    {
        return new VideoDto
        {
            VideoId = item.Id,
            Name = item.Title,
            Description = $"Streaming: {item.GameName}",
            Source = "Twitch",
            PreviewUrl = item.ThumbnailUrl?.Replace("{width}", "1280").Replace("{height}", "720") ?? "",
            VideoUrl = item.ThumbnailUrl.Replace("{width}", "640").Replace("{height}", "360"),
            CreatedAt = item.StartedAt,
            IsPrivate = false,
            VideoUploadStatus = VideoStatus.Ready,
            Views = item.ViewerCount > 0 ? item.ViewerCount : 0,
            User = new UserVideoDto
            {
                UserId = item.UserId,
                Username = item.UserName ?? item.DisplayName ?? "Unknown",
                AvatarUrl = avatarUrl ?? "",
                IsFollowed = false
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