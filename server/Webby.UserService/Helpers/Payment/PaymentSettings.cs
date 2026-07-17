namespace Webby.UserService.Helpers.Payment;

public class PaymentSettings
{
   public string PaymentPublicKey { get; set; }
   public string PaymentSecretKey { get; set; }
   public string PaymentWebhookKey { get; set; }
   public string SuccessUrl { get; set; }
   public string CancelUrl { get; set; }
   
}