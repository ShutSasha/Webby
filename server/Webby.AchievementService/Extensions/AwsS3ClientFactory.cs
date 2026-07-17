using Amazon;
using Amazon.Runtime;
using Amazon.S3;
using Webby.AchievementService.Dtos.Storage;

namespace Webby.AchievementService.Extensions;

public static class AwsS3ClientFactory
{
   public static IAmazonS3 CreateS3Client(IConfiguration configuration)
   {
      var awsOptions = configuration.GetSection(nameof(AwsOptions)).Get<AwsOptions>();
      var credentials = new BasicAWSCredentials(awsOptions?.AccessKey, awsOptions?.SecretKey);
      var config = new AmazonS3Config
      {
         RegionEndpoint = RegionEndpoint.EUNorth1
      };

      return new AmazonS3Client(credentials, config);
      
   }
}