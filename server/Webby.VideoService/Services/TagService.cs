using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services;

public class TagService : ITagService
{
   private readonly ITagRepository _tagRepository;
   
   public TagService(ITagRepository tagRepository)
   {
      _tagRepository = tagRepository;
   }
   
   
}