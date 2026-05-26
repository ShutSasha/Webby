namespace Webby.AchievementService.Interfaces.Services;

public interface IStorageService
{
   Task<string> UploadFileAsync(Guid id, string folder, string key, Stream fileStream, string contentType);
   Task DeleteFileAsync(string path);
   
}