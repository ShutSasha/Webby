using System.ComponentModel;
using Microsoft.AspNetCore.Mvc;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Stream.Enums;

namespace Webby.VideoService.Dtos.Stream;

public class SearchStreamOptions : SearchOptions
{
   [FromQuery(Name="searchPlatform")] 
   public SearchStreamPlatforms? SearchPlatform { get; set; }

   [FromQuery(Name = "nextPageToken")]
   public string? NextPageToken { get; set; }
}