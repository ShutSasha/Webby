using Webby.UserService.Models.Enums;

namespace Webby.UserService.Models;

public class Payment
{
   public Guid PaymentId { get; set; }
   public Guid UserId { get; set; }
   public decimal Amount { get; set; }
   public string Currency { get; set; }
   public string? ExternalId { get; set; }
   public DateTime CreatedAt { get; set; }
   public PaymentStatus Status { get; set; }
}