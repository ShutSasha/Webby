using System.ComponentModel.DataAnnotations;
using Webby.VideoService.Helpers.Validation;

namespace Webby.VideoService.Dtos.Video;

public class CreateVideoRequest
{
   [Required] 
   public Guid VideoId { get; set; }
   
   [Required]
   public string Name { get; set; }
   
   [Required]
   public string? Description { get; set; }
   

   [MaxFileSize(20 * 1024 * 1024, ErrorMessage = "Preview size can't exceed 20MB")]
   [Required]
   public IFormFile PreviewFile { get; set; }
   
   [Required]
   public bool IsPrivate { get; set; }
   
   [MaxLength(5, ErrorMessage = "Count of video tags can't be more than 5")]
   public List<string>? VideoTags { get; set; }
}