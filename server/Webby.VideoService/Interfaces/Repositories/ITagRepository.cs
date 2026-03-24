using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Repositories;

public interface ITagRepository : IRepository<Models.Tag>
{
   Task<List<Tag>> GetTagsByNames(List<string> names);
   Task AddTags(List<Tag> tags);
   Task AddVideoTags(List<VideoTag> videoTags);
   Task<List<Tag>> GetVideoTags(List<Guid>? videoTags);
}