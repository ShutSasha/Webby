using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;

namespace Webby.VideoService.Services;

public class TagService : ITagService
{
   private readonly ITagRepository _tagRepository;
   
   public TagService(ITagRepository tagRepository)
   {
      _tagRepository = tagRepository;
   }
   
   public async Task EnsureCreateTags(List<string>? tags, Guid videoId)
   {
      if (tags == null || !tags.Any())
         return;

      var normalizedTags = tags
         .Select(t => t.Trim().ToLower())
         .Distinct()
         .ToList();
      
      var existingTags = await _tagRepository.GetTagsByNames(normalizedTags);

      var existingNames = existingTags
         .Select(t => t.Name)
         .ToHashSet();
      
      var newTags = normalizedTags
         .Where(t => !existingNames.Contains(t))
         .Select(t => new Tag
         {
            TagId = Guid.NewGuid(),
            Name = t
         })
         .ToList();

      if (newTags.Any())
         await _tagRepository.AddTags(newTags);

      var allTags = existingTags.Concat(newTags).ToList();
      
      var videoTags = allTags
         .Select(tag => new VideoTag
         {
            VideoId = videoId,
            TagId = tag.TagId
         })
         .ToList();

      await _tagRepository.AddVideoTags(videoTags);
   }

   public async Task<List<string>> GetTagNames(List<VideoTag>? videoTags)
   {
      var tags = await _tagRepository.GetVideoTags(videoTags.Select(vt => vt.TagId).ToList());
      return tags
         .Select(t => t.Name)
         .ToList();
   }

   public async Task SyncVideoTags(Video video, List<string> tagNames)
   {
      var normalized = tagNames
         .Select(t => t.Trim().ToLower())
         .Distinct()
         .ToList();

      var existingTags = await _tagRepository.GetTagsByNames(normalized);

      var existingTagNames = existingTags
         .Select(t => t.Name.ToLower())
         .ToHashSet();
      
      var newTags = normalized
         .Where(n => !existingTagNames.Contains(n))
         .Select(n => new Tag { Name = n })
         .ToList();

      if (newTags.Any())
      {
         await _tagRepository.AddTags(newTags);
         existingTags.AddRange(newTags);
      }
      
      var finalTags = existingTags;

      var finalTagIds = finalTags.Select(t => t.TagId).ToHashSet();
      var currentTagIds = video.VideoTags.Select(vt => vt.TagId).ToHashSet();
      
      var toRemove = video.VideoTags
         .Where(vt => !finalTagIds.Contains(vt.TagId))
         .ToList();

      foreach (var vt in toRemove)
      {
         video.VideoTags.Remove(vt);
      }
      
      var toAdd = finalTags
         .Where(t => !currentTagIds.Contains(t.TagId))
         .Select(t => new VideoTag
         {
            VideoId = video.VideoId,
            TagId = t.TagId
         });

      foreach (var vt in toAdd)
      {
         video.VideoTags.Add(vt);
      }
   }
}