using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Mvc;

namespace Webby.NotificationService.Dtos.Pagination;

public class PaginationRequest
{

   [DefaultValue(1)]
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   [FromQuery(Name = "page")]
   public int Page { get; set; } = 1;

   [DefaultValue(10)]
   [Range(1,double.MaxValue, ErrorMessage ="Field {0} must be greater than {1}")]
   [FromQuery(Name = "pageSize")]
   public int PageSize { get; set; } = 10;
}