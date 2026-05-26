using System.Text.Json;
using Microsoft.EntityFrameworkCore;
using StackExchange.Redis;
using Webby.AchievementService.Data;
using Webby.AchievementService.Dtos.Event;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Services.Handlers;

public class AchievementHandler
{
    private readonly IDatabase _redisDb;
    private readonly AppDbContext _dbContext;

    public AchievementHandler(IConnectionMultiplexer redis, AppDbContext dbContext)
    {
        _redisDb = redis.GetDatabase();
        _dbContext = dbContext;
    }

    public async Task HandleEventAsync(string eventType, string payloadJson)
    {
        var options = new JsonSerializerOptions { PropertyNamingPolicy = JsonNamingPolicy.CamelCase };
        var basePayload = JsonSerializer.Deserialize<BaseEventPayload>(payloadJson, options);

        if (basePayload == null || basePayload.UserId == Guid.Empty)
        {
            return;
        }

        var achievements = await _dbContext.Achievements
            .Where(a => a.EventType == eventType)
            .ToListAsync();

        if (!achievements.Any())
        {
            return;
        }

        foreach (var ach in achievements)
        {
            string redisKey = $"user:{basePayload.UserId}:progress";
            string fieldKey = ach.AchievementId.ToString();
            long currentValue = 0;

            if (basePayload.Operation == EventOperation.Absolute)
            {
                await _redisDb.HashSetAsync(redisKey, fieldKey, basePayload.Value);
                currentValue = basePayload.Value;
            }
            else
            {
                currentValue = await _redisDb.HashIncrementAsync(redisKey, fieldKey, basePayload.Value);
            }

            var progress = await _dbContext.UserAchievementProgresses
                .FirstOrDefaultAsync(p => p.UserId == basePayload.UserId && p.AchievementId == ach.AchievementId);

            if (progress == null)
            {
                progress = new UserAchievementProgress
                {
                    UserId = basePayload.UserId,
                    AchievementId = ach.AchievementId,
                    CurrentValue = (int)currentValue
                };
                _dbContext.UserAchievementProgresses.Add(progress);
            }
            else
            {
                progress.CurrentValue = (int)currentValue;
            }

            await _dbContext.SaveChangesAsync();

            if (currentValue >= ach.TargetValue)
            {
                await UnlockAchievementAsync(basePayload.UserId, ach);
            }
        }
    }

    private async Task UnlockAchievementAsync(Guid userId, Achievement achievement)
    {
        Console.WriteLine($"Achievement unlocked: {achievement.Title} for user {userId}");
    }
}