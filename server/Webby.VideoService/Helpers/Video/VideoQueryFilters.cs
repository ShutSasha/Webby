using System.Linq.Expressions;
using Webby.VideoService.Data;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Helpers.Video;


public static class VideoQueryFilters
{
   public static (string Sql, object[] Params,
      Func<AppDbContext, Expression<Func<Models.Video, bool>>> PredicateFactory
      ) ForPlaylist(Guid playlistId, Guid? requestUserId)
   {
      return (
         @"
   (
       ""VideoId"" IN (
           SELECT ""VideoId""
           FROM ""PlaylistVideos""
           WHERE ""PlaylistId"" = {2}
       )
       AND ""IsPublished"" = TRUE
       AND (
           ""UserId"" = {3}
           OR ""IsPrivate"" = FALSE
       )
       AND (
           ""VideoUploadStatus"" = {4}
           OR (""UserId"" = {3} AND ""VideoUploadStatus"" = {5})
       )
   )
   ",
         new object[] 
         { 
            playlistId, 
            requestUserId ?? (object)DBNull.Value,
            (int)VideoStatus.Ready, 
            (int)VideoStatus.Uploading 
         },
         context => v => context.PlaylistVideos
                            .Where(pv => pv.PlaylistId == playlistId)
                            .Select(pv => pv.InternalContentId)
                            .Contains(v.VideoId)
                         && (v.UserId == requestUserId || !v.IsPrivate)
                         && (
                            v.VideoUploadStatus == VideoStatus.Ready || 
                            (v.UserId == requestUserId && v.VideoUploadStatus == VideoStatus.Uploading)
                         )
                         && v.IsPublished == true
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
          ""IsPublished"" = TRUE
         ""IsBanned"" = FALSE
          AND (
              ""IsPrivate"" = FALSE
              OR ""UserId"" = {2}
          )
          AND (
              ""VideoUploadStatus"" = {3}
          )
      )
      ",
         new object[] 
         { 
            requestUserId ?? (object)DBNull.Value,
            VideoStatus.Ready.ToString()
         },
         context => v => v.IsPublished == true
                         && (!v.IsPrivate || v.UserId == requestUserId)
                         && v.VideoUploadStatus == VideoStatus.Ready
                         && !v.IsBanned
      );
   }
}
