using System.Text.Json.Serialization;
using Webby.VideoService.Dtos.User;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Models.Enums;
using Webby.VideoService.Services.Background;

namespace Webby.VideoService.Dtos.Video;

public class VideoDto
{
   public string VideoId { get; set; }
   public required string Name { get; set; }
   public string? Description { get; set; }
   
   public string Source { get; set; } = SearchPlatforms.Webby.ToString();
   public int Views { get; set; }
   public required string PreviewUrl { get; set; }
   public required bool IsPrivate { get; set; }
   public DateTime CreatedAt { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   
   [JsonConverter(typeof(TimeSpanToStringConverter))]
   public TimeSpan Duration { get; set; }
   public string? VideoUrl { get; set; }
   public List<string>? VideoTags { get; set; }
   public UserVideoDto? User { get; set; }
}