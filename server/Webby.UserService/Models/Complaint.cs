using Webby.UserService.Models.Enums;

namespace Webby.UserService.Models;

public class Complaint
{
   public Guid ComplaintId { get; set; }
   public Guid AuthorId { get; set; }
   public User Author { get; set; }
   public ComplaintTargetType TargetType { get; set; }
   public Guid TargetId { get; set; }
   public string ReasonType { get; set; } = null!;
   public string? AdditionalInfo { get; set; }
   public DateTime CreatedAt { get; set; }
}