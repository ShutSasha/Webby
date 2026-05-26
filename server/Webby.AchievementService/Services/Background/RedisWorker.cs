using StackExchange.Redis;
using Webby.AchievementService.Services.Handlers;

namespace Webby.AchievementService.Services.Background;

public class RedisWorker : BackgroundService
{
   private readonly IConnectionMultiplexer _redis;
   private readonly IServiceScopeFactory _scopeFactory;
   private const string StreamName = "events:platform";
   private const string GroupName = "achievement-service-group";
   private const string ConsumerName = "worker-1";

   public RedisWorker(IConnectionMultiplexer redis, IServiceScopeFactory scopeFactory)
   {
      _redis = redis;
      _scopeFactory = scopeFactory;
   }

   protected override async Task ExecuteAsync(CancellationToken stoppingToken)
   {
      var db = _redis.GetDatabase();

      try { await db.StreamCreateConsumerGroupAsync(StreamName, GroupName, StreamPosition.NewMessages); }
      catch (RedisServerException ex) when (ex.Message.Contains("BUSYGROUP")) { }

      while (!stoppingToken.IsCancellationRequested)
      {
         var messages = await db.StreamReadGroupAsync(StreamName, GroupName, ConsumerName, ">", count: 10);

         foreach (var msg in messages)
         {
            var type = msg.Values.FirstOrDefault(v => v.Name == "type").Value.ToString();
            var payload = msg.Values.FirstOrDefault(v => v.Name == "payload").Value.ToString();

            using (var scope = _scopeFactory.CreateScope())
            {
               var handler = scope.ServiceProvider.GetRequiredService<AchievementHandler>();
               await handler.HandleEventAsync(type, payload);
            }

            await db.StreamAcknowledgeAsync(StreamName, GroupName, msg.Id);
         }

         if (messages.Length == 0)
         {
            await Task.Delay(500, stoppingToken);
         }
      }
   }
}