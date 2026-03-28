using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task<UploadVideoResponse> UploadVideoFile(Guid userId, UploadVideoRequest request);
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, Guid videoId);
   Task<GetVideoInformationResponse> GetVideoInformation(Guid videoId, Guid? userId);
   Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId,Guid? requestUserId,int page,int pageSize);
   Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request);
   Task<PagedResponse<VideoDto>> SearchVideo(Guid? requestUserId, SearchOptions options);
   Task<PagedResponse<VideoDto>> SearchVideoInPlaylist(Guid? requestUserId, Guid playlistId, SearchOptions searchOptions);
   Task<bool> CheckPrivateVideos(List<Guid> videoIds, Guid requestUserId);
   Task CancelVideoUploading(Guid requestUserId, Guid videoId);
   Task<bool> CheckUploadStatus(Guid videoId);
}