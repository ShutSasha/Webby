using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Dtos.Video.Enums;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task<Video> GetVideoById(Guid videoId);
   Task<UploadVideoResponse> UploadVideoFile(Guid userId, UploadVideoRequest request);
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, string videoId);
   Task<VideoDto> GetVideoInformation(string videoId, Guid? userId);
   Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId, Guid? requestedUserId, GetUserVideosRequest request);
   Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request);
   Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchVideoOptions options);
   Task<PagedResponse<VideoDto>> SearchVideoInPlaylist(Guid? requestUserId, Guid playlistId, SearchOptions searchOptions); 
   Task<PagedResponse<PreviewVideoDto>> GetRecommendationVideos(string videoId, Guid? requestUserId, int contentSeed,int page = 1, int size = 20);
   Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId);
   Task CancelVideoUploading(Guid requestUserId, string videoId);
   Task<bool> CheckUploadStatus(string videoId);
   Task IncrementVideoView(Guid requestUserId, string videoId);
   Task<(List<VideoDto> videoDtos,List<string> unavailableVideos)> GetVideoRange(List<string> ids);
}