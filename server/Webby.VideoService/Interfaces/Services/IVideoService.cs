using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Helpers.Response;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, Guid videoId);
   Task<GetVideoInformationResponse> GetVideoInformation(Guid videoId, Guid? userId);

   Task<PagedResponse<VideoDto>> GetUserVideos(Guid userId,Guid? requestUserId,int page,int pageSize);
   Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request);

}