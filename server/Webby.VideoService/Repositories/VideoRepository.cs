using Microsoft.EntityFrameworkCore;
using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;

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
         .Where(v => v.UserId == userId && (isOwner || !v.IsPrivate));

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
    ";

      var parameters = new List<Npgsql.NpgsqlParameter>
      {
         new Npgsql.NpgsqlParameter("@playlistId", playlistId),
         new Npgsql.NpgsqlParameter("@requestUserId", (object?)requestUserId ?? DBNull.Value)
      };

      if (!string.IsNullOrWhiteSpace(searchText) && searchText.Trim().Length >= 3)
      {
         query += " AND (v.\"Name\" <% @search OR v.\"Name\" ILIKE @likePattern) ";
         parameters.Add(new Npgsql.NpgsqlParameter("@search", searchText));
         parameters.Add(new Npgsql.NpgsqlParameter("@likePattern", $"%{searchText}%"));
      }
      
      var countQuery = $"SELECT COUNT(*) FROM ({query}) AS w";
      var total = await _context.Videos.FromSqlRaw(countQuery, parameters.ToArray()).CountAsync();
      
      query += " ORDER BY " +
               (!string.IsNullOrWhiteSpace(searchText) && searchText.Trim().Length >= 3
                  ? "(CASE WHEN v.\"Name\" ILIKE @likePattern THEN 1 ELSE 0 END) DESC, " +
                    "word_similarity(@search, v.\"Name\") DESC, "
                  : "") +
               "pv.\"CreatedAt\" DESC " +
               "LIMIT @take OFFSET @skip";

      parameters.Add(new Npgsql.NpgsqlParameter("@take", take));
      parameters.Add(new Npgsql.NpgsqlParameter("@skip", skip));

      var items = await _context.Videos
         .FromSqlRaw(query, parameters.ToArray())
         .ToListAsync();

      return (items, total);
   }
}