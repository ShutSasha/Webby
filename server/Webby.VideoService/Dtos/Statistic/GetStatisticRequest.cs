using Microsoft.AspNetCore.Mvc;
using Swashbuckle.AspNetCore.Annotations;

namespace Webby.VideoService.Dtos.Statistic;

public class GetStatisticRequest
{
   [FromQuery] 
   [SwaggerSchema("Used for monthly statistics. The year to filter by (e.g., 2026). Defaults to the current year if not provided.")]
   public int? Year { get; set; }
   
   [FromQuery]
   [SwaggerSchema("Used for monthly statistics. The month to filter by (1-12). Defaults to the current month if not provided.")]
   public int? Month { get; set; }
   
   [FromQuery] 
   [SwaggerSchema("Used for hourly statistics. The day of the month (1-31). If provided, the API returns hourly statistics for this specific day.")]
   public int? Day { get; set; }
}