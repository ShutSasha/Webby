using Microsoft.AspNetCore.Mvc;
using Stripe;
using Stripe.Terminal;
using Swashbuckle.AspNetCore.Annotations;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Jwt;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Controllers;

[ApiController]
[Route("api/payments")]
public class PaymentController : ControllerBase
{
   private readonly IPaymentService _paymentService;

   public PaymentController(IPaymentService paymentService)
   {
      _paymentService = paymentService;
   }

   [HttpGet("subscriptions/{paymentId:guid}")]
   [SwaggerOperation("Get user premium details after redirecting on success payment url","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<GetUserPremiumInformationResponse>>> GetPaymentDetails(Guid paymentId)
   {
      var requestedUserId = JwtHelper.ExtractUserId(HttpContext)!;
      var userPremiumInformation = await _paymentService.GetUserPremiumInformation(paymentId, requestedUserId.Value);
      return Ok(ApiResponse<GetUserPremiumInformationResponse>.Ok("Successfully retrieve premium information",
         userPremiumInformation));
   }

   [HttpPost]
   [SwaggerOperation("Initialize payment form and gives payment form url","AUTH REQUIRED")]
   public async Task<ActionResult<ApiResponse<string>>> CreatePayment()
   {
      var requestUserId = JwtHelper.ExtractUserId(HttpContext)!;
      var sessionUrl = await _paymentService.CreateCheckoutSession(requestUserId.Value);
      return Ok(ApiResponse<string>.Ok("Successfully create payment",sessionUrl));
   }
   
   
   [SwaggerIgnore]
   [HttpPost("callback")]
   public async Task<IActionResult> Handle()
   {
      var json = await new StreamReader(HttpContext.Request.Body).ReadToEndAsync();
      var signatureHeader = Request.Headers["Stripe-Signature"];
      try
      {
         await _paymentService.ProcessWebhook(json, signatureHeader!);
         return Ok();
      }
      catch (StripeException)
      {
         return BadRequest();
      }
      catch (Exception)
      {
         return StatusCode(500);
      }
   }
   
}