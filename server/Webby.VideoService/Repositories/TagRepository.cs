using Microsoft.EntityFrameworkCore;
using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Models;

namespace Webby.VideoService.Repositories;

public class TagRepository : GenericRepository<Models.Tag>, ITagRepository
{
   public TagRepository(AppDbContext context) : base(context)
   {
   }
   
   public async Task<List<Tag>> GetTagsByNames(List<string> names)
   {
      return await _context.Tags
         .Where(t => names.Contains(t.Name))
         .ToListAsync();
   }

   public async Task AddTags(List<Tag> tags)
   {
      await _context.Tags.AddRangeAsync(tags);
      await _context.SaveChangesAsync();
   }

   public async Task AddVideoTags(List<VideoTag> videoTags)
   {
      await _context.VideoTags.AddRangeAsync(videoTags);
      await _context.SaveChangesAsync();
   }

   public async Task<List<Tag>> GetVideoTags(List<Guid>? videoTags)
   {
      if (videoTags == null || !videoTags.Any())
         return [];

      return await _context.Tags
         .AsNoTracking()
         .Where(t => videoTags.Contains(t.TagId))
         .ToListAsync();
   }
}