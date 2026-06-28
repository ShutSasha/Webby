using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.VideoGrpcServer;

namespace Webby.VideoService.Services.Grpc;

public class VideoGrpcService : VideoGrpcServer.VideoGrpcService.VideoGrpcServiceBase
{
   private readonly IVideoRepository _videoRepository;
   private readonly IVideoService _videoService;

   public VideoGrpcService(IVideoRepository videoRepository, IVideoService videoService)
   {
       _videoRepository = videoRepository;
       _videoService = videoService;
   }

   public override async Task<BoolValue> CheckVideoExists(CheckVideoExistRequest request, ServerCallContext context)
   {
       if (!Guid.TryParse(request.RequestUserId, out var requestUserIdGuid))
           throw new ApiException("Check video exist error", 400, "Invalid type of request user id");

       var parseResult = PlatformPrefixToPlatformConverter.ParseLocalPlatform(request.VideoId)
                         ?? throw new ApiException("Check video exist error", 400, "Invalid video prefix type");
    
       if (!Guid.TryParse(parseResult.ActualId, out var videoIdGuid))
           throw new ApiException("Check video exist error", 400, "Invalid actual video id format");

       var video = await _videoRepository.FindById(videoIdGuid);

       if (video == null)
           return new BoolValue { Value = false };

       if (video.UserId == requestUserIdGuid)
           throw new ApiException("Check video exist error", 400, "You can't leave complaint on your own video");
   
       return new BoolValue { Value = true };
   }

   public override async Task<Empty> BanVideo(BanVideoRequest request, ServerCallContext context)
   {
       await _videoService.BanVideo(request.VideoID);
       return new Empty();
   }

   public override async Task<GetUserStatisticResponse> GetUserStatistic(GetUserStatisticRequest request, ServerCallContext context)
   {
       if (!Guid.TryParse(request.UserId, out var userId))
       {
           throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid UserId format. Must be a valid GUID."));
       }

       var now = DateTime.UtcNow;
       var targetYear = request.HasYear && request.Year > 0 ? request.Year : now.Year;
       var targetMonth = request.HasMonth && request.Month is >= 1 and <= 12 ? request.Month : now.Month;
       var daysInMonth = DateTime.DaysInMonth(targetYear, targetMonth);
       
       var totalWatched = await _videoRepository.CountUserWatchedVideosAsync(userId);
       var totalTime = await _videoRepository.GetUserTotalWatchTimeAsync(userId);
       var topTags = await _videoRepository.GetUserTopTagsAsync(userId, 5);
       var dbViewsTrend = await _videoRepository.GetUserDailyWatchTrendForMonthAsync(userId, targetYear, targetMonth);

       var response = new GetUserStatisticResponse
       {
           TotalWatchedVideos = totalWatched,
           TotalWatchTime = totalTime
       };

       response.TopTags.AddRange(topTags.Select(t => new TagStatistic
       {
           TagName = t.TagName,
           WatchCount = t.WatchCount
       }));

       var trendDict = dbViewsTrend.ToDictionary(v => v.Date, v => v.ViewsCount);
       
       var fullMonthViews = Enumerable.Range(1, daysInMonth)
           .Select(d => new DateTime(targetYear, targetMonth, d, 0, 0, 0, DateTimeKind.Utc))
           .Select(date => new DailyWatchActivity
           {
               Date = Timestamp.FromDateTime(date), 
               ViewsCount = trendDict.GetValueOrDefault(date, 0)
           });

       response.WatchActivityTrend.AddRange(fullMonthViews);

       return response;
   }
   
}