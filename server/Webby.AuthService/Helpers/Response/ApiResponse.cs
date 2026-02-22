namespace Webby.AuthService.Helpers.Response;

public class ApiResponse
{
   public bool Success { get; init; }
   public string Message { get; init; } = string.Empty;
   public Dictionary<string, string>? Errors { get; init; }

   public static ApiResponse Ok(string message)
      => new()
      {
         Success = true,
         Message = message
      };

   public static ApiResponse Fail(string message, Dictionary<string, string>? errors = null)
      => new()
      {
         Success = false,
         Message = message,
         Errors = errors
      };
}

public class ApiResponse<T> : ApiResponse
{
   public T? Data { get; init; }

   public static ApiResponse<T> Ok(string message, T data)
      => new()
      {
         Success = true,
         Message = message,
         Data = data
      };
}
