using Grpc.Core;
using Grpc.Core.Interceptors;
using Webby.VideoService.Helpers.Exception;

namespace Webby.VideoService.Middlewares;

public class GrpcClientExceptionInterceptor : Interceptor
{
   public override AsyncUnaryCall<TResponse> AsyncUnaryCall<TRequest, TResponse>(
      TRequest request,
      ClientInterceptorContext<TRequest, TResponse> context,
      AsyncUnaryCallContinuation<TRequest, TResponse> continuation)
   {
      var call = continuation(request, context);
      
      return new AsyncUnaryCall<TResponse>(
         HandleResponse(call.ResponseAsync),
         call.ResponseHeadersAsync,
         call.GetStatus,
         call.GetTrailers,
         call.Dispose);
   }

   private async Task<TResponse> HandleResponse<TResponse>(Task<TResponse> inner)
   {
      try
      {
         return await inner;
      }
      catch (RpcException ex)
      {
         var statusCode = ex.StatusCode switch
         {
            StatusCode.NotFound => 404,
            StatusCode.InvalidArgument => 400,
            StatusCode.ResourceExhausted => 400,
            StatusCode.FailedPrecondition => 400,
            StatusCode.Unauthenticated => 401,
            StatusCode.PermissionDenied => 403,
            _ => 500
         };

         throw new ApiException("gRPC Client Error", statusCode, ex.Status.Detail);
      }
   }
}