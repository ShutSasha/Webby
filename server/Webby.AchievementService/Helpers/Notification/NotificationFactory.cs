using Webby.AchievementService.Dtos.Notification;
using Webby.AchievementService.Interfaces.Helpers.Notification;

namespace Webby.AchievementService.Helpers.Notification;

public class NotificationFactory : INotificationFactory
{
   
   public Task<SendNotificationDto> UnlockAchievementSendMessage(string achievementName, Guid userId, Guid achievementId)
   {
      return Task.FromResult(new SendNotificationDto
      {
         Title = "New achievement unlocked!",
         Message = $"You unlocked {achievementName} achievement",
         UserId = userId.ToString(),
         TargetIdentifier = $"achievement-{achievementId}"
      });
   }
}