namespace Webby.AuthService.Helpers.Exception;

public class ApiException : System.Exception
{
   public int StatusCode { get; }

   public ApiException(string message, int statusCode) : base(message)
   {
      StatusCode = statusCode;
   }
}