using Microsoft.AspNetCore.Mvc;
using Stripe;
using Stripe.Terminal;
using Swashbuckle.AspNetCore.Annotations;
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