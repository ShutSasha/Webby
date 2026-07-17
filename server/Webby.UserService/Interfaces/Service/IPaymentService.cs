using Webby.UserService.Dtos.User;

namespace Webby.UserService.Interfaces.Service;

public interface IPaymentService
{
   Task<string> CreateCheckoutSession(Guid userId);
   Task ProcessWebhook(string json, string signature);
   Task<GetUserPremiumInformationResponse> GetUserPremiumInformation(Guid paymentId, Guid requestUserId);

}