using Microsoft.AspNetCore.Mvc;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Video.Enums;

namespace Webby.VideoService.Dtos.Video;

public class SearchVideoOptions : SearchOptions
{
   [FromQuery(Name = "searchPlatform")]
   public SearchVideoPlatforms? SearchPlatform { get; set; }
   
   [FromQuery(Name = "nextPageToken")]
   public string? NextPageToken { get; set; }
}