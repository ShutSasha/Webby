using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Mvc;

namespace Webby.VideoService.Dtos.Search;

public class SearchOptions
{
   [FromQuery(Name = "searchText")]
   public string? SearchText { get; set; }
   
   [DefaultValue(1)]
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   [FromQuery(Name = "page")]
   public int Page { get; set; } = 1;
   
   [DefaultValue(10)]
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   [FromQuery(Name = "pageSize")]
   public int PageSize { get; set; } = 10;

   [FromQuery(Name = "contentSeed")] 
   public int ContentSeed { get; set; } = 0;
}