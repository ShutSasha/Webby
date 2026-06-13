namespace Webby.VideoService.Dtos.Statistic;

public class UserStatisticDto
{
   public int TotalWatchedVideos { get; set; }
   public long TotalWatchTime { get; set; } 
   public IReadOnlyCollection<TagStatisticDto> TopTags { get; set; } = Array.Empty<TagStatisticDto>();
   public IReadOnlyCollection<DailyViewsDto> WatchActivityTrend { get; set; } = Array.Empty<DailyViewsDto>();
}

public class TagStatisticDto
{
   public string TagName { get; set; } = string.Empty;
   public int WatchCount { get; set; }

   public TagStatisticDto() { }

   public TagStatisticDto(string tagName, int watchCount)
   {
      TagName = tagName;
      WatchCount = watchCount;
   }
}