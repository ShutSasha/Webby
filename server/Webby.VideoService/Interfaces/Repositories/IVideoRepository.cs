using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IVideoRepository : IRepository<Video>
{
   Task<Video> GetVideoInformationById(Guid videoId);
   Task<(List<Video>, int)> GetPaginatedUserVideos(Guid userId,bool isOwner,int page,int pageSize);
   Task<(List<Video> Items, int Total)> SearchAsync(string? searchText, int skip, int take);
}