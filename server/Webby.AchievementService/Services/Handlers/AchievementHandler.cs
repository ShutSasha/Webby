using System.Text.Json;
using Microsoft.EntityFrameworkCore;
using StackExchange.Redis;
using Webby.AchievementService.Data;
using Webby.AchievementService.Dtos.Event;
using Webby.AchievementService.Interfaces.Helpers.Notification;
using Webby.AchievementService.Interfaces.Services;
using Webby.AchievementService.Models;
using Webby.NotificationService.GrpcClient;

namespace Webby.AchievementService.Services.Handlers;

public class AchievementHandler
{
    private readonly IDatabase _redisDb;
    private readonly AppDbContext _dbContext;
    private readonly NotificationGrpcService.NotificationGrpcServiceClient _notificationGrpcServiceClient;
    private readonly INotificationFactory _notificationFactory;
    private readonly IEventTypeService _eventTypeService;
    
    public AchievementHandler(IConnectionMultiplexer redis, AppDbContext dbContext, NotificationGrpcService.NotificationGrpcServiceClient notificationGrpcServiceClient, INotificationFactory notificationFactory, IEventTypeService eventTypeService)
    {
        _redisDb = redis.GetDatabase();
        _dbContext = dbContext;
        _notificationGrpcServiceClient = notificationGrpcServiceClient;
        _notificationFactory = notificationFactory;
        _eventTypeService = eventTypeService;
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

        if (achievements.Count == 0)
        {
            return;
        }

        foreach (var ach in achievements)
        {
            var alreadyUnlocked = await _dbContext.UserAchievements
                .AnyAsync(ua => ua.UserId == basePayload.UserId && ua.AchievementId == ach.AchievementId);

            if (alreadyUnlocked)
            {
                continue;
            }

            var progress = await _dbContext.UserAchievementProgresses
                .FirstOrDefaultAsync(p => p.UserId == basePayload.UserId && p.AchievementId == ach.AchievementId);

            if (!basePayload.IsIncrementOperation && progress != null && progress.CurrentValue >= basePayload.Value)
            {
                continue;
            }
                
            var redisKey = $"user:{basePayload.UserId}:progress";
            var fieldKey = ach.AchievementId.ToString();
            long currentValue = 0;

            if (basePayload.IsIncrementOperation)
            {
                currentValue = await _redisDb.HashIncrementAsync(redisKey, fieldKey, basePayload.Value);
            }
            else
            {
                await _redisDb.HashSetAsync(redisKey, fieldKey, basePayload.Value);
                currentValue = basePayload.Value;
            }

            if (progress == null)
            {
                progress = new UserAchievementProgress
                {
                    UserId = basePayload.UserId,
                    AchievementId = ach.AchievementId,
                    CurrentValue = (int)currentValue
                };
               await _dbContext.UserAchievementProgresses.AddAsync(progress);
            }
            else
            {
                progress.CurrentValue = (int)currentValue;
            }

            await _dbContext.SaveChangesAsync();

            if (currentValue >= ach.TargetValue)
            {
                await UnlockAchievementAsync(basePayload.UserId, ach);
                await _redisDb.HashDeleteAsync(redisKey, fieldKey);
                _dbContext.UserAchievementProgresses.Remove(progress);
            }
            
            if (!await _eventTypeService.HasEventType(eventType))
                await _eventTypeService.AddEventType(eventType);
        }
    }
    

    private async Task UnlockAchievementAsync(Guid userId, Achievement achievement)
    {
        var userAchievement = new UserAchievement()
        {
            AchievementId = achievement.AchievementId,
            IsPinned = false,
            UnlockedAt = DateTime.UtcNow,
            UserId = userId,
        };

        _dbContext.UserAchievements.Add(userAchievement);
        await _dbContext.SaveChangesAsync();
        
        var notificationDto = await _notificationFactory
            .UnlockAchievementSendMessage(achievement.Title, userId, achievement.AchievementId);
        
        await _notificationGrpcServiceClient.SendNotificationToUserAsync(new CreateNotificationRequest
        {
            Message = notificationDto.Message,
            TargetType = GrpcNotificationTargetType.Achievement,
            TargetIdentifier = notificationDto.TargetIdentifier,
            Title = notificationDto.Title,
            UserId = notificationDto.UserId
        });
        
    }
}