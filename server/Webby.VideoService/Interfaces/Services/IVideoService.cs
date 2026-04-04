using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task<Video> GetVideoById(Guid videoId);
   Task<UploadVideoResponse> UploadVideoFile(Guid userId, UploadVideoRequest request);
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, Guid videoId);
   Task<VideoDto> GetVideoInformation(string videoId, Guid? userId, SearchPlatforms platform);
   Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId,Guid? requestUserId,int page,int pageSize);
   Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request);
   Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchVideoOptions options);
   Task<PagedResponse<VideoDto>> SearchVideoInPlaylist(Guid? requestUserId, Guid playlistId, SearchOptions searchOptions);
   Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId);
   Task CancelVideoUploading(Guid requestUserId, Guid videoId);
   Task<bool> CheckUploadStatus(Guid videoId);
   Task<GlobalSearchVideoResponse> GlobalSearchVideos(Guid? requestUserId, GlobalSearchOptions searchOptions);
   Task IncrementVideoView(Guid requestUserId, Guid videoId);
}