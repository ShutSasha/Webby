using Grpc.Core;
using Webby.MediaService.GrpcServer;
using Webby.VideoService.Constants;
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

   public MediaGrpcService(IVideoService videoService,ITwitchSearchService twitchSearchService,
      IYouTubeSearchService youtubeSearchService)
   {
      _videoService = videoService;
      _twitchSearchService = twitchSearchService;
      _youtubeSearchService = youtubeSearchService;
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
            
            return new VideoResponse()
            {
               Id = stream.StreamerId,
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

   public override async Task<GetAuthorIDByVideoIDResponse> GetAuthorIDByVideoID(GetAuthorIDByVideoIDRequest request, ServerCallContext context)
   {
      var platformResult = PlatformPrefixToPlatformConverter.ParseLocalPlatform(request.VideoID);

      if (platformResult == null)
      {
         throw new ApiException("Get author error", 400, "Invalid platform prefix type");
      }
      
      
      var (_, actualId) = platformResult.Value;
      
      if (!Guid.TryParse(actualId, out var videoIdGuid))
      {
         throw new ApiException("Get author error", 400, "Invalid video id type");
      }
      
      var videoResult = await _videoService.GetVideoById(videoIdGuid,showPrivate: true);
      
      return new GetAuthorIDByVideoIDResponse() { AuthorID = videoResult.UserId.ToString() };
   }
}