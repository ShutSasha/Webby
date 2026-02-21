namespace Webby.UserService.Helpers.Exception;

public class ApiException : System.Exception
{
   public int StatusCode { get; }
   public Dictionary<string, string>? Errors { get; }

   public ApiException(string message, int statusCode,
      Dictionary<string, string>? errors = null) : base(message)
   {
      StatusCode = statusCode;
      Errors = errors;
   }
}