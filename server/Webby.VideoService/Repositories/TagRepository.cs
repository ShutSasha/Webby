using Webby.VideoService.Data;
using Webby.VideoService.Interfaces.Repositories;

namespace Webby.VideoService.Repositories;

public class TagRepository : GenericRepository<Models.Tag>, ITagRepository
{
   public TagRepository(AppDbContext context) : base(context)
   {
   }
}