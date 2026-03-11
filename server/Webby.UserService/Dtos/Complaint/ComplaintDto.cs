namespace Webby.UserService.Dtos.Complaint;

public class ComplaintDto
{
   public Guid ComplaintId { get; set; }
   public string ReasonType { get; set; }
   public string AdditionalInfo { get; set; }
   public DateTime CreatedAt { get; set; }
}