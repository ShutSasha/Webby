using System.Text.Json;
using Microsoft.Extensions.Options;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.External;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services;

public class YoutubeSearchService : IExternalVideoSearchService
{
    private readonly HttpClient _httpClient;
    private readonly string _apiKey;

    public YoutubeSearchService(
        HttpClient httpClient,
        IOptions<ExternalServicesOptions> options)
    {
        _httpClient = httpClient;
        _apiKey = options.Value.YouTubeOptions.ApiKey;
    }

    public async Task<PagedResponse<VideoDto>> SearchAsync(SearchVideoOptions options)
{
    var url = BuildUrl(options);
    var response = await _httpClient.GetAsync(url);
    response.EnsureSuccessStatusCode();

    var initialData = await response.Content.ReadFromJsonAsync<YouTubeVideoResponse>(
        new JsonSerializerOptions { PropertyNameCaseInsensitive = true });

    if (initialData?.Items == null || !initialData.Items.Any())
        return new PagedResponse<VideoDto>();
    
    var videoIds = initialData.Items.Select(item => 
    {
        if (item.Id is JsonElement element && element.ValueKind == JsonValueKind.Object)
            return element.TryGetProperty("videoId", out var idProp) ? idProp.GetString() : null;
        
        return item.Id?.ToString();
    }).Where(id => !string.IsNullOrEmpty(id)).ToList();
    
    var detailsUrl = DefaultLinks.BaseYouTubeVideosLink +
                     $"?part=snippet,contentDetails,statistics" +
                     $"&id={string.Join(",", videoIds)}" +
                     $"&key={_apiKey}";

    var detailsResponse = await _httpClient.GetAsync(detailsUrl);
    detailsResponse.EnsureSuccessStatusCode();

    var jsonDetailsResponse = await detailsResponse.Content.ReadAsStringAsync();

    var fullData = await detailsResponse.Content.ReadFromJsonAsync<YouTubeVideoResponse>(
        new JsonSerializerOptions { PropertyNameCaseInsensitive = true });
    
    var videoDtos = fullData.Items.Select(item => new VideoDto
    {
        VideoId = item.Id?.ToString(),
        Name = item.Snippet.Title,
        Description = item.Snippet.Description,
        Source = "YouTube",
        PreviewUrl = item.Snippet.Thumbnails?.Medium?.Url ?? "",
        VideoUrl = $"https://www.youtube.com/watch?v={item.Id}",
        CreatedAt = item.Snippet.PublishedAt,
        IsPrivate = false,
        VideoUploadStatus = VideoStatus.Ready,

        Duration = !string.IsNullOrEmpty(item.ContentDetails?.Duration) 
            ? System.Xml.XmlConvert.ToTimeSpan(item.ContentDetails.Duration) 
            : TimeSpan.Zero,
            
        Views = int.TryParse(item.Statistics?.ViewCount, out var v) ? v : 0,
        VideoTags = item.Snippet.Tags,
        
        User = new UserVideoDto
        {
            UserId = item.Snippet.ChannelId,
            Username = item.Snippet.ChannelTitle,
            AvatarUrl = DefaultLinks.BaseYouTubeUserIcon,
            IsFollowed = false
        }
    }).ToList();

    return new PagedResponse<VideoDto>
    {
        Items = videoDtos,
        NextPageToken = initialData.NextPageToken,
        PageSize = initialData.PageInfo.ResultsPerPage,
        TotalCount = initialData.PageInfo.TotalResults
    };
}

    private string BuildUrl(SearchVideoOptions options)
    {
        var hasSearch = !string.IsNullOrWhiteSpace(options.SearchText);
        
        string part = "id"; 
        string url;

        if (hasSearch)
        {
            var query = Uri.EscapeDataString(options.SearchText);
            url = DefaultLinks.BaseYouTubeSearchLink +
                  $"?part={part}" +
                  $"&type=video" +
                  $"&q={query}" +
                  $"&maxResults={options.PageSize}" +
                  $"&key={_apiKey}";
        }
        else
        {
            url = DefaultLinks.BaseYouTubeVideosLink +
                  $"?part={part}" +
                  "&chart=mostPopular" +
                  "&regionCode=US" + 
                  $"&maxResults={options.PageSize}" +
                  $"&key={_apiKey}";
        }

        if (!string.IsNullOrEmpty(options.NextPageToken))
        {
            url += $"&pageToken={options.NextPageToken}";
        }

        return url;
    }
}