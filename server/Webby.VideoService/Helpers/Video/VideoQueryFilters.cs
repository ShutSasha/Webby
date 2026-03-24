using System.Linq.Expressions;
using Webby.VideoService.Data;

namespace Webby.VideoService.Helpers.Video;


public static class VideoQueryFilters
{
   public static (string Sql, object[] Params,
      Func<AppDbContext, Expression<Func<Models.Video, bool>>> PredicateFactory
      ) ForPlaylist(Guid playlistId,Guid? requestUserId)
   {
      return (
         @"
         (
             ""VideoId"" IN (
                 SELECT ""VideoId""
                 FROM ""PlaylistVideos""
                 WHERE ""PlaylistId"" = {2}
             )
             AND (
                 ""UserId"" = {3}
                 OR ""IsPrivate"" = FALSE
             )
         )
         ",
         new object[] { playlistId, requestUserId },
         context => v => context.PlaylistVideos
                            .Where(pv => pv.PlaylistId == playlistId)
                            .Select(pv => pv.VideoId)
                            .Contains(v.VideoId)
                         && (v.UserId == requestUserId || !v.IsPrivate)
      );
   }

   public static (
      string Additional,
      object[] Params,
      Func<AppDbContext, Expression<Func<Models.Video, bool>>> PredicateFactory
      ) SearchVideoFilter(Guid? requestUserId)
   {
      return (
         @"
        (
            ""IsPrivate"" = FALSE
            OR ""UserId"" = {2}
        )
        ",
         new object[] { requestUserId },
         context => v => !v.IsPrivate || v.UserId == requestUserId
      );
   }
}
