using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;

namespace Webby.UserService.Dtos.Statistic;

public class GetStatisticRequest
{
   [FromQuery]
   [SwaggerSchema("Used for monthly statistics. The year to filter by (e.g., 2026). Defaults to current year.")]
   public int? Year { get; set; }

   [FromQuery]
   [SwaggerSchema("Used for monthly statistics. The month to filter by (1-12). Defaults to current month.")]
   public int? Month { get; set; }
}

public class UserStatisticDto
{
   public int TotalWatchedVideos { get; set; }
   public long TotalWatchTime { get; set; }
   public IReadOnlyCollection<TagStatisticDto> TopTags { get; set; } = Array.Empty<TagStatisticDto>();
   public IReadOnlyCollection<DailyViewsDto> WatchActivityTrend { get; set; } = Array.Empty<DailyViewsDto>();
}

public record TagStatisticDto(string TagName, int WatchCount);
public record DailyViewsDto(DateTime Date, int ViewsCount);