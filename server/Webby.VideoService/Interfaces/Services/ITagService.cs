namespace Webby.VideoService.Interfaces.Services;

public interface ITagService
{
   Task EnsureCreateTags(List<string>? tags, Guid videoId);
}