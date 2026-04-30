using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IVideoRepository : IRepository<Video>
{
   Task<Video> GetVideoInformationById(Guid videoId);
   Task<(List<Video>, int)> GetPaginatedUserVideos(Guid userId,bool isOwner, int page, int pageSize, bool shouldShowDrafts);
   Task<bool> CheckVideosCount(List<Guid> videoIds);
   Task<bool> CheckForbiddenVideos(List<Guid> playlistVideosIds, Guid requestUserId);
   Task<(List<Video> Items, int Total)> SearchVideosInPlaylistAsync(
      Guid playlistId,
      Guid? requestUserId,
      string? searchText,
      int skip,
      int take);

   Task UpdateVideoFileMetaData(Guid videoId, string videoFileUrl, long duration);
   Task<bool> FindUserView(Guid userId, Guid videoId);
   Task AddUserView(UserView userView);
   Task<int> CountUserView(Guid videoId);
   Task<(List<Video> Items, int Total)> GetRecommendedVideosAsync(Guid? currentVideoId, List<string> currentTags,
      List<Guid> subscribedIds, List<string> historyTags, int skip, int pageSize);

   Task<List<string>> GetRecentUserViewTagsAsync(Guid userId, int limit = 30);

}