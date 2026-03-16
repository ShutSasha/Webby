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
}