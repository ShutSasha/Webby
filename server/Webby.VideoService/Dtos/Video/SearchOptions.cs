using System.ComponentModel.DataAnnotations;

namespace Webby.VideoService.Dtos.Video;

public class SearchOptions
{
   public string? SearchText { get; set; }
   
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   public int Page { get; set; } = 1;
   
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   public int PageSize { get; set; } = 10;
}