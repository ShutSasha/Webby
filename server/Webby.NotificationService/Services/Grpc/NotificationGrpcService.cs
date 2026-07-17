using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Microsoft.AspNetCore.SignalR;
using Webby.NotificationService.Constants;
using Webby.NotificationService.GrpcService;
using Webby.NotificationService.Helpers.Exception;
using Webby.NotificationService.Hubs;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Interfaces.Services;
using Webby.NotificationService.Models;
using Webby.NotificationService.Models.Enums;

namespace Webby.NotificationService.Services.Grpc;

public class NotificationGrpcService: GrpcService.NotificationGrpcService.NotificationGrpcServiceBase
{
   private readonly INotificationRepository _notificationRepository;
   private readonly IHubContext<NotificationHub> _hubContext;

   public NotificationGrpcService(INotificationRepository notificationRepository, IHubContext<NotificationHub> hubContext)
   {
      _notificationRepository = notificationRepository;
      _hubContext = hubContext;
   }
   public override async Task<CreateNotificationResponse> SendNotificationToUser(CreateNotificationRequest request, ServerCallContext context)
   {
      if (string.IsNullOrWhiteSpace(request.UserId) || !Guid.TryParse(request.UserId, out var userId))
      {
         throw new ApiException("Send notification error", 400, "Invalid format of user id");
      }
      
      if (request.TargetType == GrpcNotificationTargetType.Unspecified)
      {
         throw new ApiException("Send notification error",400,"Target type can't be unspecified");
      }
      
      var domainTargetType = (NotificationTargetType)request.TargetType;
      
      var isDuplicate = await _notificationRepository
         .IsRecentDuplicateAsync(userId, domainTargetType, request.TargetIdentifier, TimeSpan.FromHours(2));
      
      if (isDuplicate)
      {
         return new CreateNotificationResponse { Success = true };
      }
      
      var notification = new Notification
      {
         NotificationId = Guid.NewGuid(),
         CreatedAt = DateTime.UtcNow,
         Message = request.Message,
         NotificationStatus = NotificationStatus.Unread,
         TargetIdentifier = request.TargetIdentifier,
         TargetType = domainTargetType,
         Title = request.Title,
         UserId = userId
      };

      try
      {
         await _notificationRepository.Add(notification);
         await _hubContext.Clients.User(request.UserId).SendAsync(WebSocketMethodNames.ReceiveNotificationMethod, notification);

         var unreadNotificationsCount = await _notificationRepository.CountNotifications(userId, NotificationStatus.Unread);
         await _hubContext.Clients.User(request.UserId).SendAsync(WebSocketMethodNames.UpdateUnreadMessagesMethod, unreadNotificationsCount);
      }
      catch (Exception)
      {
         throw new ApiException("Send notification error", 500, "Something went wrong");
      }
      
      return new CreateNotificationResponse
      {
         Success = true
      };
   }

   public override async Task<Empty> ReportBlocking(ReportBlockingRequest request, ServerCallContext context)
   {
      await _hubContext.Clients.User(request.UserId).SendAsync(WebSocketMethodNames.UserBanSystemNotificationMethod, request.UserId);

      return new Empty();
   }
}