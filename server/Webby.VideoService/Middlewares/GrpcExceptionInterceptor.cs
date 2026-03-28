using Grpc.Core;
using Grpc.Core.Interceptors;
using Webby.VideoService.Helpers.Exception;

namespace Webby.VideoService.Middlewares;

public class GrpcExceptionInterceptor : Interceptor
{
   public override async Task<TResponse> UnaryServerHandler<TRequest, TResponse>(
      TRequest request,
      ServerCallContext context,
      UnaryServerMethod<TRequest, TResponse> continuation)
   {
      try
      {
         return await continuation(request, context);
      }
      catch (ApiException apiEx)
      {
         var grpcStatus = apiEx.StatusCode switch
         {
            400 => StatusCode.InvalidArgument,
            401 => StatusCode.Unauthenticated,
            403 => StatusCode.PermissionDenied,
            404 => StatusCode.NotFound,
            _ => StatusCode.Internal
         };
         
         var detail = apiEx.Errors != null 
            ? string.Join("; ", apiEx.Errors.Select(x => $"{x.Key}: {x.Value}"))
            : apiEx.Message;

         throw new RpcException(new Status(grpcStatus, detail));
      }
      catch (Exception ex)
      {
         throw new RpcException(new Status(StatusCode.Internal, $"An internal error occurred {ex.Message}"));
      }
   }
}