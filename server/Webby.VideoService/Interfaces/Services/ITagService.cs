using Webby.VideoService.Models;

namespace Webby.VideoService.Interfaces.Services;

public interface ITagService
{
   Task EnsureCreateTags(List<string>? tags, Guid videoId);
   Task<List<string>> GetTagNames(List<VideoTag>? videoTags);
   Task SyncVideoTags(Video video, List<string> tagNames);

}