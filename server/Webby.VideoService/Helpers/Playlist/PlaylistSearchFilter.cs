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
                ""IsPrivate"" = FALSE
                OR ""UserId"" = {2}
            )
            ",
         new object[] { userId },
         context => p => !p.IsPrivate || p.UserId == userId
      );
   }
}