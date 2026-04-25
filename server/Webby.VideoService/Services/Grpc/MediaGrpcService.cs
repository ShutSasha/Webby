using Grpc.Core;
using Webby.MediaService.GrpcServer;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services.Grpc;

public class MediaGrpcService : MediaService.GrpcServer.MediaService.MediaServiceBase
{
   private readonly IVideoService _videoService;
   private readonly IPlaylistService _playlistService;
   private readonly YoutubeSearchService _youtubeSearchService;
   private readonly TwitchSearchService _twitchSearchService;

   public MediaGrpcService(IVideoService videoService, IPlaylistService playlistService, TwitchSearchService twitchSearchService, YoutubeSearchService youtubeSearchService)
   {
      _videoService = videoService;
      _playlistService = playlistService;
      _twitchSearchService = twitchSearchService;
      _youtubeSearchService = youtubeSearchService;
   }

   public override async Task<VideoResponse> GetVideo(GetVideoRequest request, ServerCallContext context)
{
    switch (request.VideoType)
    {
        case VideoType.Video:
        {
            if (!Guid.TryParse(request.Id, out var videoId))
            {
                throw new ApiException("Validation error", 400, "Invalid ID format");
            }

            var video = await _videoService.GetVideoById(videoId);
            
            return new VideoResponse
            {
                Id = video.VideoId.ToString(),
                Title = video.Name,
                Thumbnail = video.PreviewUrl,
                VideoUrl = video.VideoUrl
            };
        }
        
        case VideoType.Youtube:
        {
            var ytVideo = await _youtubeSearchService.FindById(request.Id);
            
            return new VideoResponse
            {
                Id = ytVideo.VideoId, 
                Title = ytVideo.Name,
                Thumbnail = ytVideo.PreviewUrl,
                VideoUrl = ytVideo.VideoUrl 
            };
        }
        
        case VideoType.Twitch:
        {
            var stream = await _twitchSearchService.FindById(request.Id);
            
                return new VideoResponse
                {
                    Id = stream.StreamId,
                    Title = stream.Name,
                    Thumbnail = stream.PreviewUrl,
                    VideoUrl = stream.StreamUrl
                };
                
        }
        default:
            throw new ApiException("Validation error", 400, "Unsupported video type");
    }
}

   public override async Task<PlaylistResponse> GetPlaylist(GetPlaylistRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.Id, out var playlistId))
      {
         throw new ApiException("Validation error", 400, "Invalid ID format");
      }
      
      var playlistInfo = await _playlistService.GetPlaylistInformation(playlistId, null);
      
      var videosResult = await _videoService.SearchVideoInPlaylist(null,playlistId, new SearchOptions
      {
         Page = request.Page,
         PageSize = request.PageSize
      });
      
      var response = new PlaylistResponse
      {
         Id = playlistInfo.Playlist.PlaylistId.ToString(),
         Title = playlistInfo.Playlist.Name,
         Thumbnail = playlistInfo.FirstVideo?.PreviewUrl ?? DefaultLinks.PlaylistEmptyLink,
         TotalCount = videosResult.TotalCount
      };
      
      var videoNodes = videosResult.Items.Select(v => new VideoResponse
      {
         Id = v.VideoId.ToString(),
         Title = v.Name,
         Thumbnail = v.PreviewUrl,
         VideoUrl = v.VideoUrl
      });

      response.Videos.AddRange(videoNodes);

      return response;
   }
   
}