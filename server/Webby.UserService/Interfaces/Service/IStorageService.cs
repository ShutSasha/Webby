namespace Webby.UserService.Interfaces.Service;

public interface IStorageService
{
   Task<string> UploadFileAsync(Guid id, string key, Stream fileStream, string contentType);
   Task DeleteFileAsync(Guid id, string path);
   
}