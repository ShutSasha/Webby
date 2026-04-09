using Microsoft.Extensions.Options;
using Stripe;
using Stripe.Checkout;
using Webby.UserService.Consts;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Helpers.Payment;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;
using Webby.UserService.Models.Enums;

namespace Webby.UserService.Services;

public class PaymentService : IPaymentService
{
    private readonly IPaymentRepository _paymentRepository;
    private readonly IUserPremiumRepository _userPremiumRepository;
    private readonly PaymentSettings _settings;

    public PaymentService(IPaymentRepository paymentRepository, IUserPremiumRepository userPremiumRepository, IOptions<PaymentSettings> paymentSettings)
    {
        _paymentRepository = paymentRepository;
        _userPremiumRepository = userPremiumRepository;
        _settings = paymentSettings.Value;
        StripeConfiguration.ApiKey = _settings.PaymentSecretKey;
    }

    public async Task<string> CreateCheckoutSession(Guid userId)
    {
        var payment = new Payment
        {
            PaymentId = Guid.NewGuid(),
            UserId = userId,
            Amount = 1m,
            Currency = "EUR",
            Status = PaymentStatus.Pending,
            CreatedAt = DateTime.UtcNow
        };

        await _paymentRepository.Add(payment);

        var options = new SessionCreateOptions
        {
            SuccessUrl = _settings.SuccessUrl + $"/{payment.PaymentId}",
            CancelUrl = _settings.CancelUrl,
            PaymentMethodTypes = ["card"],
            LineItems =
            [
                new SessionLineItemOptions
                {
                    PriceData = new SessionLineItemPriceDataOptions
                    {
                        UnitAmount = (long)(payment.Amount * 100),
                        Currency = payment.Currency,
                        ProductData = new SessionLineItemPriceDataProductDataOptions
                        {
                            Name = "Webby premium status",
                            Description = "Premium status that allows more customization feature and unlimited using",
                            Images =[DefaultLinks.PaymentImageLink]
                        },
                    },
                    Quantity = 1,
                }

            ],
            Mode = "payment",
            ClientReferenceId = payment.PaymentId.ToString()
        };

        var service = new SessionService();
        var session = await service.CreateAsync(options);

        payment.ExternalId = session.Id;
        await _paymentRepository.Update(payment);

        return session.Url;
    }

    public async Task ProcessWebhook(string json, string signature)
    {
        try
        {
            var stripeEvent =
                EventUtility.ConstructEvent(json, signature, _settings.PaymentWebhookKey, throwOnApiVersionMismatch: false); //падает на этой строчке

            switch (stripeEvent.Type)
            {
                case EventTypes.CheckoutSessionCompleted:
                    var session = stripeEvent.Data.Object as Session;
                    await HandleSuccess(session);
                    break;

                case EventTypes.CheckoutSessionExpired:
                    var expiredSession = stripeEvent.Data.Object as Session;
                    await HandleFailure(expiredSession.Id, PaymentStatus.Canceled);
                    break;

                case EventTypes.PaymentIntentPaymentFailed:
                    var intent = stripeEvent.Data.Object as PaymentIntent;
                    await HandleFailure(intent.Id, PaymentStatus.Failed);
                    break;
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }
       
    }

    public async Task<GetUserPremiumInformationResponse> GetUserPremiumInformation(Guid paymentId, Guid requestUserId)
    {
        var payment = await _paymentRepository.FindById(paymentId)
                      ?? throw new ApiException("Get payment information error",404,"Payment wasn't found");

        if (payment.UserId != requestUserId)
        {
            throw new ApiException("Get user information error", 403, "You can't get information of another user");
        }

        var userPremiumInformation = await _userPremiumRepository.GetUserPremiumInformation(payment.UserId);

        if (userPremiumInformation == null)
        {
            throw new ApiException("Get information error",404,"User premium wasn't found");
        }

        return new GetUserPremiumInformationResponse
        {
            ExpirationDate = userPremiumInformation.ExpiresAt,
            Username = userPremiumInformation.User.Username
        };
    }

    private async Task HandleSuccess(Session session)
    {
        var payment = await _paymentRepository.GetByExternalId(session.Id);
        if (payment.Status != PaymentStatus.Pending) return;
        
        payment.Status = PaymentStatus.Succeeded;
        payment.ExternalId = session.PaymentIntentId;
        await _paymentRepository.Update(payment);
        
        var premium = await _userPremiumRepository.FindById(payment.UserId);
        if (premium == null)
        {
            await _userPremiumRepository.Add(new UserPremium {
                UserId = payment.UserId,
                ExpiresAt = DateTime.UtcNow.AddMonths(1)
            });
        }
        else
        {
            var start = premium.ExpiresAt > DateTime.UtcNow ? premium.ExpiresAt : DateTime.UtcNow;
            premium.ExpiresAt = start.AddMonths(1);
            await _userPremiumRepository.Update(premium);
        }
    }

    private async Task HandleFailure(string externalId, PaymentStatus status)
    {
        var payment = await _paymentRepository.GetByExternalId(externalId);
        if (payment is { Status: PaymentStatus.Pending })
        {
            payment.Status = status;
            await _paymentRepository.Update(payment);
        }
    }
    
}