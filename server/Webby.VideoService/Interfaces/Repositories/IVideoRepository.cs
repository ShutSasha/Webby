using Webby.VideoService.Dtos.Statistic;
using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface IVideoRepository : IRepository<Video>
{
   Task<Video> GetVideoInformationById(Guid videoId);
   Task<int> CountUserVideos(Guid userId);
   
   Task<(List<Video>, int)> GetPaginatedUserVideos(Guid userId,bool isOwner, int page, int pageSize, bool shouldShowDrafts);
   Task<bool> CheckVideosCount(List<Guid> videoIds);
   Task<bool> CheckForbiddenVideos(List<Guid> playlistVideosIds, Guid requestUserId);
   Task<(List<PlaylistVideo> Items, int Total)> SearchVideosInPlaylistAsync(
      Guid playlistId,
      Guid? requestUserId,
      string? searchText,
      int skip,
      int take);
   Task UpdateVideoFileMetaData(Guid videoId, string videoFileUrl, long duration);
   Task<bool> FindUserView(Guid userId, Guid videoId);
   Task AddUserView(UserView userView);
   Task<int> CountUserView(Guid videoId);
   Task<(List<Video> Items, int Total, int Seed)> GetRecommendedVideosAsync(Guid? currentVideoId, List<string> currentTags,
      List<Guid> subscribedIds, List<string> historyTags, int skip, int pageSize, int contentSeed = 0);
   Task<List<string>> GetRecentUserViewTagsAsync(Guid userId, int limit = 30);
   Task<List<DailyViewsDto>> GetVideoDailyViewsTrendForMonthAsync(Guid videoId, int year, int month);
   Task<List<HourlyActivityDto>> GetVideoHourlyActivityForDayAsync(Guid videoId, int year, int month, int day);
   Task<int> CountUserWatchedVideosAsync(Guid userId);
   Task<long> GetUserTotalWatchTimeAsync(Guid userId);
   Task<List<TagStatisticDto>> GetUserTopTagsAsync(Guid userId, int limit = 5);
   Task<List<DailyViewsDto>> GetUserDailyWatchTrendForMonthAsync(Guid userId, int year, int month);

}