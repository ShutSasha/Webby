using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;
using Webby.VideoService.Dtos.Playlist;

namespace Webby.VideoService.Dtos.Video;

public class SearchVideoOptions : SearchOptions
{
   [FromQuery(Name = "searchPlatform")]
   public SearchPlatforms SearchPlatform { get; set; }
   
   [FromQuery(Name = "nextPageToken")]
   public string? NextPageToken { get; set; }
}