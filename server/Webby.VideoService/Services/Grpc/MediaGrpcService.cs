using Grpc.Core;
using Webby.MediaService.GrpcServer;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Platforms.Enums;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Helpers.Converters;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Services.Grpc;

public class MediaGrpcService : MediaService.GrpcServer.MediaService.MediaServiceBase
{
   private readonly IVideoService _videoService;
   private readonly IYouTubeSearchService _youtubeSearchService;
   private readonly ITwitchSearchService _twitchSearchService;
   private readonly ILogger<MediaGrpcService> _logger;

   public MediaGrpcService(IVideoService videoService,ITwitchSearchService twitchSearchService,
      IYouTubeSearchService youtubeSearchService, ILogger<MediaGrpcService> logger)
   {
      _videoService = videoService;
      _twitchSearchService = twitchSearchService;
      _youtubeSearchService = youtubeSearchService;
      _logger = logger;
   }


   public override async Task<VideoResponse> GetVideo(GetVideoRequest request, ServerCallContext context)
   {
      var parseResult = PlatformPrefixToPlatformConverter.ParseSystemPlatform(request.Id);

      if (parseResult == null)
      {
         throw new ApiException("Get video error", 400, "Invalid id format");
      }

      var (videoPlatform, actualId) = parseResult.Value;

      switch (videoPlatform)
      {
         case SystemPlatforms.Webby:
         {
            if (!Guid.TryParse(actualId, out var videoIdGuid))
            {
               throw new ApiException("Get video error", 400, "Invalid webby video id format");
            }

            var video = await _videoService.GetVideoById(videoIdGuid);

            return new VideoResponse()
            {
               Id = PlatformPrefixesConstants.WebbyPrefix + video.VideoId,
               Thumbnail = video.PreviewUrl,
               Title = video.Name,
               VideoUrl = video.VideoUrl
            };
            
         }
         case SystemPlatforms.YouTube:
         {
            var ytVideo = await _youtubeSearchService.FindById(actualId);

            return new VideoResponse()
            {
               Id = ytVideo.VideoId,
               Thumbnail = ytVideo.PreviewUrl,
               Title = ytVideo.Name,
               VideoUrl = ytVideo.VideoUrl
            };
         }
         case SystemPlatforms.Twitch:
         {
            var stream = await _twitchSearchService.FindById(actualId);
            
            _logger.LogInformation("Stream Id {id}, stream thumbnail {thumbnail} steam title {title} url {url}",
               stream.StreamId, stream.PreviewUrl,stream.Name,stream.StreamUrl);
            
            return new VideoResponse()
            {
               Id = stream.StreamId,
               Thumbnail = stream.PreviewUrl,
               Title = stream.Name,
               VideoUrl = stream.StreamUrl
            };
         }
         
         default:
            throw new ApiException("Get video error", 400, "Invalid video platform type");
      }
      
   }

   public override async Task<GetVideosBatchResponse> GetVideosBatch(GetVideosBatchRequest request, ServerCallContext context)
   {
      var getRangeResult = await _videoService.GetVideoRange(request.Ids.ToList());

      var response = new GetVideosBatchResponse();

      var mappedVideos = getRangeResult.videoDtos.Select(v => new VideoResponse
      {
         Id = v.VideoId ?? string.Empty,
         Title = v.Name ?? string.Empty,
         Thumbnail = v.PreviewUrl ?? string.Empty,
         VideoUrl = v.VideoUrl ?? string.Empty
      });
      
      response.Videos.AddRange(mappedVideos);
      response.UnavailableVideoIds.AddRange(getRangeResult.unavailableVideos);

      return response;
   }
}