using System.Linq.Expressions;
using Webby.VideoService.Data;

namespace Webby.VideoService.Helpers.Video;


public static class VideoQueryFilters
{
   public static (
      string Sql,
      object[] Params,
      Func<AppDbContext, Expression<Func<Models.Video, bool>>> PredicateFactory
      ) ForPlaylist(Guid playlistId)
   {
      return (
         @"
         (
             ""VideoId"" IN (
                 SELECT ""VideoId""
                 FROM ""PlaylistVideos""
                 WHERE ""PlaylistId"" = {2}
             )
             OR ""IsPrivate"" = FALSE
         )
         ",
         new object[] { playlistId },
         context => v => context.PlaylistVideos
                            .Where(pv => pv.PlaylistId == playlistId)
                            .Select(pv => pv.VideoId)
                            .Contains(v.VideoId)
                         && !v.IsPrivate
      );
   }
}
