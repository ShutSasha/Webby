using Webby.AchievementService.Dtos.Notification;

namespace Webby.AchievementService.Interfaces.Helpers.Notification;

public interface INotificationFactory
{
   Task<SendNotificationDto> UnlockAchievementSendMessage(string achievementName, Guid userId, Guid achievementId);
}