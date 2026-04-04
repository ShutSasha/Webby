using Microsoft.AspNetCore.Components;
using Microsoft.AspNetCore.Mvc;
using Webby.VideoService.Dtos.Playlist;

namespace Webby.VideoService.Dtos.Search;

public class GlobalSearchOptions : SearchOptions
{
   [FromQuery(Name = "sectionType")]
   public SearchSections SectionType { get; set; }
}