using System.Text.Json.Serialization;

namespace Webby.VideoService.Dtos.External;

public class TwitchResponse<T>
{
   [JsonPropertyName("data")]
   public List<T> Data { get; set; } = new();

   [JsonPropertyName("pagination")]
   public TwitchPagination? Pagination { get; set; }
}

public class TwitchPagination
{
   [JsonPropertyName("cursor")]
   public string? Cursor { get; set; }
}

public class TwitchItem
{
   [JsonPropertyName("id")] public string Id { get; set; } = string.Empty;
   [JsonPropertyName("user_id")] public string? UserId { get; set; }
   [JsonPropertyName("user_name")] public string? UserName { get; set; }
   [JsonPropertyName("game_name")] public string? GameName { get; set; }
   [JsonPropertyName("viewer_count")] public int ViewerCount { get; set; }
   
   [JsonPropertyName("display_name")] public string? DisplayName { get; set; }
   [JsonPropertyName("broadcaster_login")] public string? BroadcasterLogin { get; set; }
   [JsonPropertyName("title")] public string Title { get; set; } = string.Empty;
   [JsonPropertyName("thumbnail_url")] public string ThumbnailUrl { get; set; } = string.Empty;
   [JsonPropertyName("started_at")] public DateTime StartedAt { get; set; }
}

public class TwitchUser
{
   [JsonPropertyName("id")] public string Id { get; set; } = string.Empty;
   [JsonPropertyName("profile_image_url")] public string ProfileImageUrl { get; set; } = string.Empty;
   [JsonPropertyName("display_name")] public string DisplayName { get; set; } = string.Empty;
}
