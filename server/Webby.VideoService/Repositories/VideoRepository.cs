using Microsoft.EntityFrameworkCore;
using Npgsql;
using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Repositories;

public class VideoRepository : GenericRepository<Video>,IVideoRepository
{
   public VideoRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<Video> GetVideoInformationById(Guid videoId)
   {
      return await _context.Videos
         .Include(v => v.VideoTags)
         .FirstAsync(v => v.VideoId == videoId);
   }
   public async Task<(List<Video>, int)> GetPaginatedUserVideos(Guid userId, bool isOwner, int page, int pageSize)
   {
      var query = _context.Videos
         .Where(v => v.UserId == userId 
                     && (isOwner || !v.IsPrivate)
                     && (v.VideoUploadStatus == VideoStatus.Ready || (isOwner && v.VideoUploadStatus == VideoStatus.Uploading)));

      var totalCount = await query.CountAsync();

      var videos = await query
         .OrderByDescending(v => v.CreatedAt)
         .Skip((page - 1) * pageSize)
         .Take(pageSize)
         .ToListAsync();

      return (videos, totalCount);
   }

   public async Task<bool> CheckVideosCount(List<Guid> videoIds)
   {
      return await _context.Videos.Where(v => videoIds.Contains(v.VideoId)).CountAsync() == videoIds.Count;
   }

   public async Task<bool> CheckForbiddenVideos(List<Guid> playlistVideosIds, Guid requestUserId)
   {
      return await _context.Videos.AnyAsync(v => playlistVideosIds.Contains(v.VideoId)
                                      && v.IsPrivate
                                      && v.UserId != requestUserId);
   }
   
public async Task<(List<Video> Items, int Total)> SearchVideosInPlaylistAsync(
   Guid playlistId,
   Guid? requestUserId,
   string? searchText,
   int skip,
   int take)
{
   var query = @"
     SELECT v.*
     FROM ""PlaylistVideos"" pv
     JOIN ""Videos"" v ON v.""VideoId"" = pv.""VideoId""
     WHERE pv.""PlaylistId"" = @playlistId
       AND (v.""IsPrivate"" = FALSE OR v.""UserId"" = @requestUserId)
       AND (
           v.""VideoUploadStatus"" = @statusReady 
           OR (v.""UserId"" = @requestUserId AND v.""VideoUploadStatus"" = @statusUploading)
       )
 ";
   
   var parameters = new List<NpgsqlParameter>
   {
      new NpgsqlParameter("@playlistId", playlistId),
      new NpgsqlParameter("@requestUserId", (object?)requestUserId ?? DBNull.Value),
      new NpgsqlParameter("@statusReady", VideoStatus.Ready.ToString()),
      new NpgsqlParameter("@statusUploading", VideoStatus.Uploading.ToString())
   };

   if (!string.IsNullOrWhiteSpace(searchText) && searchText.Trim().Length >= 3)
   {
      query += " AND (v.\"Name\" <% @search OR v.\"Name\" ILIKE @likePattern) ";
      parameters.Add(new NpgsqlParameter("@search", searchText));
      parameters.Add(new NpgsqlParameter("@likePattern", $"%{searchText}%"));
   }
   
   var countParameters = parameters.Select(p => p.Clone()).ToArray();
   var total = await _context.Videos
      .FromSqlRaw(query, countParameters)
      .CountAsync();
   
   query += " ORDER BY " +
            (!string.IsNullOrWhiteSpace(searchText) && searchText.Trim().Length >= 3
               ? "(CASE WHEN v.\"Name\" ILIKE @likePattern THEN 1 ELSE 0 END) DESC, " +
                 "word_similarity(@search, v.\"Name\") DESC, "
               : "") +
            "pv.\"CreatedAt\" DESC " +
            "LIMIT @take OFFSET @skip";
   
   parameters.Add(new NpgsqlParameter("@take", take));
   parameters.Add(new NpgsqlParameter("@skip", skip));
   
   var items = await _context.Videos
      .FromSqlRaw(query, parameters.ToArray())
      .ToListAsync();

   return (items, total);
}
}