namespace Webby.UserService.Interfaces.Service;

public interface IPaymentService
{
   Task<string> CreateCheckoutSession(Guid userId);
   Task ProcessWebhook(string json, string signature);
   
}