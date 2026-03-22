using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.VideoGrpcServer;

namespace Webby.VideoService.Services.Grpc;

public class VideoGrpcService : VideoGrpcServer.VideoGrpcService.VideoGrpcServiceBase
{
   private readonly IVideoRepository _videoRepository;

   public VideoGrpcService(IVideoRepository videoRepository)
   {
      _videoRepository = videoRepository;
   }

   public override async Task<BoolValue> CheckVideoExists(CheckVideoExistRequest request, ServerCallContext context)
   {
      var video = await _videoRepository.FindById(Guid.Parse(request.VideoId));

      return new BoolValue { Value = video != null };
   }
}