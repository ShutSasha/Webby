using System.Text.Json;
using StackExchange.Redis;
using Webby.VideoService.Constants;
using Webby.VideoService.Interfaces.Helpers;

public class EventPublisher : IEventPublisher
{
   private readonly IDatabase _redisDb;
   private readonly ILogger<EventPublisher> _logger;


   public EventPublisher(
      IConnectionMultiplexer redis, 
      ILogger<EventPublisher> logger)
   {
      _redisDb = redis.GetDatabase();
      _logger = logger;
   }

   public async Task PublishAsync<T>(T @event) where T : IPlatformEvent
   {
      try
      {
         var payloadJson = JsonSerializer.Serialize(@event, new JsonSerializerOptions 
         { 
            PropertyNamingPolicy = JsonNamingPolicy.CamelCase 
         });
         
         var streamEntries = new[]
         {
            new NameValueEntry("type", @event.EventType),
            new NameValueEntry("payload", payloadJson),
            new NameValueEntry("timestamp", DateTimeOffset.UtcNow.ToUnixTimeMilliseconds())
         };
         
         await _redisDb.StreamAddAsync(RedisConstants.StreamName, streamEntries);
         
      }
      catch (Exception ex)
      {
         _logger.LogError(ex, "Error sending event with type: {EventType}", @event.EventType);
         throw; 
      }
   }
}