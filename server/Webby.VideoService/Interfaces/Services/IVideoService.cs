using Webby.VideoService.Dtos.Video;

namespace Webby.VideoService.Interfaces.Services;

public interface IVideoService
{
   Task CreateVideo(Guid userId, CreateVideoRequest request);
   Task DeleteVideo(Guid userId, Guid videoId);
   
}