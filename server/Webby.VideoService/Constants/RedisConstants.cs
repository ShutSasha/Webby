namespace Webby.VideoService.Constants;

public static class RedisConstants
{
   public const string StreamName = "events:platform";
   public const string UserStreamName = "events:users";
   public const int RedisConnectionTimeout = 10;
}