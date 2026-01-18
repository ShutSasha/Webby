namespace Webby.AuthService.Helpers.Response;

public class ApiResponse
{
   public bool Success { get; set; }
   public string Message { get; set; } = string.Empty;
   public object? Data { get; set; }
   public Dictionary<string, string>? Errors { get; set; }

   public static ApiResponse Ok(string message, object? data = null)
   {
      return new ApiResponse { Success = true, Message = message, Data = data };
   }

   public static ApiResponse Fail(string message, Dictionary<string, string>? errors = null)
   {
      return new ApiResponse { Success = false, Message = message, Errors = errors };
   }
}
