using Microsoft.EntityFrameworkCore;
using Npgsql;
using Webby.VideoService.Constants;
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
         .ThenInclude(vt => vt.Tag)
         .FirstAsync(v => v.VideoId == videoId);
   }

   public async Task<int> CountUserVideos(Guid userId)
   {
       return await _context.Videos
           .Where(v => v.UserId == userId && !v.IsPrivate)
           .CountAsync();
   }

   public async Task<(List<Video>, int)> GetPaginatedUserVideos(
      Guid userId, 
      bool isOwner, 
      int page, 
      int pageSize, 
      bool shouldShowDrafts)
   {
      var query = _context.Videos.Where(v => v.UserId == userId);
      
      if (isOwner)
      {
         if (!shouldShowDrafts)
         {
            query = query.Where(v => v.IsPublished && v.VideoUploadStatus == VideoStatus.Ready);
         }
      }
      else
      {
         query = query.Where(v => 
            v.IsPublished && 
            !v.IsPrivate && 
            v.VideoUploadStatus == VideoStatus.Ready);
      }
      
      var totalCount = await query.CountAsync();
      
      var videos = await query
         .OrderByDescending(v => v.CreatedAt)
         .Skip((page - 1) * pageSize)
         .Take(pageSize)
         .Include(v => v.VideoTags)!
         .ThenInclude(vt => vt.Tag)
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
   
    public async Task<(List<PlaylistVideo> Items, int Total)> SearchVideosInPlaylistAsync(
     Guid playlistId,
     Guid? requestUserId,
     string? searchText,
     int skip,
     int take)
    {
        var query = @"
         SELECT pv.*
         FROM ""PlaylistVideos"" pv
         LEFT JOIN ""Videos"" v ON v.""VideoId"" = pv.""InternalContentId""
         WHERE pv.""PlaylistId"" = @playlistId
           AND (
               pv.""Platform"" = @youtubePlatform
               OR pv.""Platform"" = @twitchPlatform
               OR (
                   pv.""Platform"" = @webbyPlatform
                   AND v.""IsPublished"" = TRUE
                   AND (v.""IsPrivate"" = FALSE OR v.""UserId"" = @requestUserId)
                   AND (
                       v.""VideoUploadStatus"" = @statusReady 
                       OR (v.""UserId"" = @requestUserId AND v.""VideoUploadStatus"" = @statusUploading)
                   )
               )
           )
     ";
       
        var parameters = new List<NpgsqlParameter>
        {
            new NpgsqlParameter("@playlistId", playlistId),
            new NpgsqlParameter("@requestUserId", (object?)requestUserId ?? DBNull.Value),
            new NpgsqlParameter("@statusReady", VideoStatus.Ready.ToString()),
            new NpgsqlParameter("@statusUploading", VideoStatus.Uploading.ToString()),
            new NpgsqlParameter("@youtubePlatform", SystemPlatforms.YouTube.ToString()),
            new NpgsqlParameter("@twitchPlatform", SystemPlatforms.Twitch.ToString()),
            new NpgsqlParameter("@webbyPlatform", SystemPlatforms.Webby.ToString())
        };
       
        var trimmedSearch = searchText?.Trim();
        bool hasSearch = !string.IsNullOrWhiteSpace(trimmedSearch);

        if (hasSearch)
        {
            if (trimmedSearch!.Length >= 3)
            {
                query += " AND (pv.\"Platform\" = @youtubePlatform OR pv.\"Platform\" = @twitchPlatform OR (pv.\"Platform\" = @webbyPlatform AND (v.\"Name\"::text <% @search OR v.\"Name\" ILIKE @likePattern))) ";
            }
            else
            {
                query += " AND (pv.\"Platform\" = @youtubePlatform OR pv.\"Platform\" = @twitchPlatform OR (pv.\"Platform\" = @webbyPlatform AND v.\"Name\" ILIKE @likePattern)) ";
            }
            parameters.Add(new NpgsqlParameter("@search", trimmedSearch));
            parameters.Add(new NpgsqlParameter("@likePattern", $"%{trimmedSearch}%"));
        }
       
        var countParameters = parameters.Select(p => p.Clone()).ToArray();
        var total = await _context.PlaylistVideos
            .FromSqlRaw(query, countParameters)
            .CountAsync();
        
        if (hasSearch && trimmedSearch!.Length >= 3)
        {
            query += " ORDER BY " +
                     "(CASE WHEN pv.\"Platform\" = @webbyPlatform AND v.\"Name\" ILIKE @likePattern THEN 1 ELSE 0 END) DESC, " +
                     "(CASE WHEN pv.\"Platform\" = @webbyPlatform THEN word_similarity(@search, v.\"Name\"::text) ELSE 0 END) DESC, " +
                     "pv.\"CreatedAt\" DESC ";
        }
        else
        {
            query += " ORDER BY pv.\"CreatedAt\" DESC ";
        }
       
        query += "LIMIT @take OFFSET @skip";
       
        parameters.Add(new NpgsqlParameter("@take", take));
        parameters.Add(new NpgsqlParameter("@skip", skip));
       
        var items = await _context.PlaylistVideos
            .FromSqlRaw(query, parameters.ToArray())
            .Include(pv => pv.Video) 
            .ToListAsync();

        return (items, total);
    }

   public async Task UpdateVideoFileMetaData(Guid videoId, string videoFileUrl, long duration)
   {
      await _context.Videos
         .Where(v => v.VideoId == videoId)
         .ExecuteUpdateAsync(s => s
            .SetProperty(v => v.VideoUrl, videoFileUrl)
            .SetProperty(v => v.Duration, duration)
         );
   }
   
   public async Task<bool> FindUserView(Guid userId, Guid videoId) 
      => await _context.UserViews.AnyAsync(uv => uv.UserId == userId && uv.VideoId == videoId);
   

   public async Task AddUserView(UserView userView)
   {
      await _context.UserViews.AddAsync(userView);
      await _context.SaveChangesAsync();
   }

   public async Task<int> CountUserView(Guid videoId) 
      => await _context.UserViews.CountAsync(uv => uv.VideoId == videoId);
   
   public async Task<List<string>> GetRecentUserViewTagsAsync(Guid userId, int limit = 30)
   {
       var recentVideoIds = await _context.UserViews
           .Where(uv => uv.UserId == userId)
           .Select(uv => uv.VideoId)
           .Take(limit)
           .ToListAsync();

       if (recentVideoIds.Count == 0)
           return [];
       
       return await _context.Videos
           .Where(v => recentVideoIds.Contains(v.VideoId))
           .SelectMany(v => v.VideoTags!)
           .Select(vt => vt.Tag.Name)
           .Distinct()
           .ToListAsync();
   }

   public async Task<(List<Video> Items, int Total, int Seed)> GetRecommendedVideosAsync(
    Guid? currentVideoId,
    List<string> currentTags,
    List<Guid> subscribedIds,
    List<string> historyTags,
    int skip,
    int pageSize,
    int contentSeed = 0)
   {
       var query = _context.Videos
           .Where(v => v.IsPublished 
                    && !v.IsPrivate 
                    && v.VideoUploadStatus == VideoStatus.Ready);
       
       if (currentVideoId.HasValue && currentVideoId.Value != Guid.Empty)
       {
           query = query.Where(v => v.VideoId != currentVideoId.Value);
       }
       
       var scoredQuery = query.Select(v => new
       {
           Video = v,
           
           CurrentTagsScore = currentTags.Any() 
               ? v.VideoTags!.Count(vt => currentTags.Contains(vt.Tag.Name)) * ScoreConstants.VideoTagWatchScoreMultiplier
               : 0,
               
           SubscriptionScore = subscribedIds.Contains(v.UserId) 
              ? ScoreConstants.UserFollowScore
              : 0,
           
           HistoryTagsScore = historyTags.Any()
               ? v.VideoTags!.Count(vt => historyTags.Contains(vt.Tag.Name)) * ScoreConstants.HistoryVideoTagScoreMultiplier
               : 0,
           
           PopularityScore = v.Views / ScoreConstants.PopularityScoreDivider 
       });
       
       var orderedQuery = scoredQuery
           .OrderByDescending(x => x.CurrentTagsScore + x.SubscriptionScore + x.HistoryTagsScore + x.PopularityScore)
           .ThenByDescending(x => x.Video.CreatedAt);
       
       var total = await orderedQuery.CountAsync();
       var seed = contentSeed == 0 ? Guid.NewGuid().GetHashCode() : contentSeed;
    
       const int maxCandidatePoolSize = 30; 
    
       var pagedIds = new List<Guid>();
       
       if (skip < maxCandidatePoolSize)
       {
           var poolIds = await orderedQuery
               .Take(maxCandidatePoolSize)
               .Select(x => x.Video.VideoId)
               .ToListAsync();

           var random = new Random(seed);
           var randomizedPool = poolIds.OrderBy(id => random.Next()).ToList();
           
           var takeFromPool = Math.Min(pageSize, maxCandidatePoolSize - skip);
           pagedIds.AddRange(randomizedPool.Skip(skip).Take(takeFromPool));
       }
       
       if (pagedIds.Count < pageSize)
       {
           var skipFromDb = Math.Max(maxCandidatePoolSize, skip); 
           var takeFromDb = pageSize - pagedIds.Count;

           var regularIds = await orderedQuery
               .Skip(skipFromDb)
               .Take(takeFromDb)
               .Select(x => x.Video.VideoId)
               .ToListAsync();

           pagedIds.AddRange(regularIds);
       }
       
       List<Video> items = [];
       if (pagedIds.Count > 0)
       {
           var unorderedItems = await _context.Videos
               .Include(v => v.VideoTags)!
               .ThenInclude(vt => vt.Tag)
               .Where(v => pagedIds.Contains(v.VideoId))
               .ToListAsync();
        
           items = unorderedItems.OrderBy(v => pagedIds.IndexOf(v.VideoId)).ToList();
       }

       return (items, total, seed);
   }
}