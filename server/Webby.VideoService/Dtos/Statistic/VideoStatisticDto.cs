using Amazon.Runtime;

namespace Webby.VideoService.Dtos.Statistic;

public class VideoStatisticDto
{
   public int TotalViews { get; set; }
   public long TotalSecondsWatched { get; set; }
   public IReadOnlyCollection<DailyViewsDto> ViewsTrend { get; set; }
   public IReadOnlyCollection<HourlyActivityDto> HourlyActivity { get; set; }
   
}

public record DailyViewsDto(DateTime Date, int ViewsCount);

public record HourlyActivityDto(int Hour, int ViewsCount);