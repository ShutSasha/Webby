using System.Linq.Expressions;
using Webby.VideoService.Data;

namespace Webby.VideoService.Helpers.Playlist;

public static class PlaylistSearchFilter
{
   public static (
      string Additional,
      object[] Params,
      Func<AppDbContext, Expression<Func<Models.Playlist, bool>>> PredicateFactory
      ) SearchPlaylistFilters(Guid? requestUserId)
   {
      return (
         @"
            (
                ""IsPrivate"" = FALSE
                OR ""UserId"" = {2}
            )
            ",
         new object[] { requestUserId },
         context => p => !p.IsPrivate || p.UserId == requestUserId
      );
   }
   
   public static (
      string Additional,
      object[] Params,
      Func<AppDbContext, Expression<Func<Models.Playlist, bool>>> PredicateFactory
      ) SearchUserPlaylistsFilter(Guid? userId, bool shouldShowPrivate)
   {
      return (
         @"
        (
            ""UserId"" = {2}
            AND (
                ""IsPrivate"" = FALSE
                OR {3} = TRUE
            )
        )
        ",
         new object[] { userId, shouldShowPrivate },
         context => p =>
            p.UserId == userId &&
            (!p.IsPrivate || shouldShowPrivate)
      );
   }
}