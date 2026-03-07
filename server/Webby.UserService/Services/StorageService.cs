using System.Net;
using Amazon.S3;
using Amazon.S3.Model;
using Microsoft.Extensions.Options;
using Webby.UserService.Dtos.Storage;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Services;

public class StorageService : IStorageService
{
   private readonly IAmazonS3 _s3Client;
   private readonly AwsOptions _options;

   public StorageService(IAmazonS3 s3Client, IOptions<AwsOptions> options)
   {
      _s3Client = s3Client;
      _options = options.Value;
   }
   
   public async Task<string> UploadFileAsync(Guid id, string key, Stream fileStream, string contentType)
   {
      var request = new PutObjectRequest
      {
         BucketName = _options.BucketName,
         Key = "user_data/" + $"{id}/" + DateTime.Now.ToString("HH:mm:ss") + key,
         InputStream = fileStream,
         ContentType = contentType,
         CannedACL = S3CannedACL.BucketOwnerFullControl
      };

      try
      {
         var response = await _s3Client.PutObjectAsync(request);
         
         if (response.HttpStatusCode == HttpStatusCode.OK)
         {
            return $"https://webby-watch-platform-bucket.s3.eu-north-1.amazonaws.com/{request.Key}";
         }
         
         throw new ApiException("File upload failed",500);
      }
      catch (AmazonS3Exception e)
      {
         throw new ApiException("S3 Upload error", 500, e.Message);
      }
      catch (Exception e)
      {
         throw new ApiException("Internal service upload error", 500, e.Message);
      }
   }
   
   public async Task DeleteFileAsync(Guid id, string path)
   {
      var uri = new Uri(path);
      var key = uri.AbsolutePath.TrimStart('/');
      
      var request = new DeleteObjectRequest
      {
         BucketName = _options.BucketName,
         Key = key 
      };

      try
      {
         var response = await _s3Client.DeleteObjectAsync(request);

         if (response.HttpStatusCode != HttpStatusCode.NoContent &&
             response.HttpStatusCode != HttpStatusCode.OK)
         {
            throw new ApiException("File deletion failed", 500);
         }
      }
      catch (AmazonS3Exception e)
      {
         throw new ApiException(e.Message, 500);
      }
      catch (Exception e)
      {
         throw new ApiException(e.Message, 500);
      }
   }
}