using Webby.VideoService.Dtos.Video;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, Guid videoId);
   Task<GetVideoInformationResponse> GetVideoInformation(Guid videoId, Guid? userId);
   Task<List<VideoDto>> GetUserVideos(Guid userId, Guid? requestUserId);
   Task UpdateVideoInformation(Guid userId, UpdateVideoRequest request);

}