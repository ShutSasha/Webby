using System.Collections;

namespace Webby.VideoService.Models;

public class Tag
{
   public Guid TagId { get; set; }
   public string Name { get; set; }
   public ICollection<VideoTag> VideoTags { get; set; }
}