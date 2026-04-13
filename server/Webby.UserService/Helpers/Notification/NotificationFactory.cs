using Webby.UserService.Dtos.Notification;
using Webby.UserService.Interfaces.Helpers;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Helpers.Notification;

public class NotificationFactory : INotificationFactory
{
   private readonly IUserRepository _userRepository;

   public NotificationFactory(IUserRepository userRepository)
   {
      _userRepository = userRepository;
   }
   public async Task<SendNotificationDto> CreateNewFollowerNotification(Guid targetUserId, Guid followerId)
   {
      var user = await _userRepository.FindById(followerId);

      return new SendNotificationDto
      {
         Title = "New follower on your profile!",
         Message = $"User {user!.Username} followed you",
         UserId = targetUserId.ToString(),
         TargetIdentifier = $"user-profile-{followerId}"
      };
   }
}