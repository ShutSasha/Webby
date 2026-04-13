using Webby.UserService.Dtos.Notification;

namespace Webby.UserService.Interfaces.Helpers;

public interface INotificationFactory
{
   Task<SendNotificationDto> CreateNewFollowerNotification(Guid targetUserId, Guid followerId);
}