using Microsoft.AspNetCore.SignalR;
using Webby.NotificationService.Constants;
using Webby.NotificationService.Dtos.Notification;
using Webby.NotificationService.Dtos.Pagination;
using Webby.NotificationService.Helpers.Exception;
using Webby.NotificationService.Helpers.Response;
using Webby.NotificationService.Hubs;
using Webby.NotificationService.Interfaces.Repositories;
using Webby.NotificationService.Interfaces.Services;
using Webby.NotificationService.Models;
using Webby.NotificationService.Models.Enums;
using Webby.NotificationService.UserGrpcClient;
using Microsoft.Extensions.Caching.Memory;

namespace Webby.NotificationService.Services;

public class NotificationService : INotificationService
{
   private readonly INotificationRepository _notificationRepository;
   private readonly IHubContext<NotificationHub> _hubContext;
   private readonly IMemoryCache _memoryCache;
   private readonly UserGrpcClient.UserGrpcService.UserGrpcServiceClient _userClient;
   public NotificationService(INotificationRepository notificationRepository,
      IHubContext<NotificationHub> hubContext, IMemoryCache memoryCache,
      UserGrpcService.UserGrpcServiceClient userClient)
   {
      _notificationRepository = notificationRepository;
      _hubContext = hubContext;
      _memoryCache = memoryCache;
      _userClient = userClient;
   }
   
   public async Task<Notification> CreateNotification(CreateNotificationRequest request)
   {
      var notification = new Notification()
      {
         NotificationId = Guid.NewGuid(),
         UserId = request.UserId,
         Title = request.Title,
         Message = request.Message,
         CreatedAt = DateTime.UtcNow,
         TargetIdentifier = request.TargetIdentifier,
         TargetType = request.TargetType,
         NotificationStatus = NotificationStatus.Unread,
      };
      await _notificationRepository.Add(notification);
      
      await _hubContext.Clients.User(request.UserId.ToString())
         .SendAsync("ReceiveNotification", notification);
      return notification;
   }

   public async Task<Notification> UpdateNotification(UpdateNotificationRequest request)
   {
      var notification = await _notificationRepository.FindById(request.NotificationId)
                         ?? throw new ApiException("Update notification error", 404, "Notification wasn't found");

      notification.Message = request.Message;
      notification.Title = request.Title;
      notification.TargetType = request.TargetType;
      notification.TargetIdentifier = request.TargetIdentifier;

      await _notificationRepository.Update(notification);

      return notification;
   }

   public async Task<PagedResponse<Notification>> GetUnreadNotifications(Guid userId, PaginationRequest request)
   {
      var skip = (request.Page - 1) * request.PageSize;

      var totalCount = await _notificationRepository.CountAsync(
         n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Unread);

      var notifications = await _notificationRepository.GetByPredicate(
         predicate: n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Unread,
         orderBy: q => q.OrderByDescending(n => n.CreatedAt),
         skip: skip,
         take: request.PageSize
      );

      return new PagedResponse<Notification>
      {
         Items = notifications.ToList(),
         TotalCount = totalCount,
         Page = request.Page,
         PageSize = request.PageSize
      };
   }

   public async Task<PagedResponse<Notification>> GetReadNotifications(Guid userId, PaginationRequest request)
   {
      var skip = (request.Page - 1) * request.PageSize;

      var totalCount = await _notificationRepository.CountAsync(
         n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Read);

      var notifications = await _notificationRepository.GetByPredicate(
         predicate: n => n.UserId == userId && n.NotificationStatus == NotificationStatus.Read,
         orderBy: q => q.OrderByDescending(n => n.CreatedAt),
         skip: skip,
         take: request.PageSize
      );

      return new PagedResponse<Notification>
      {
         Items = notifications.ToList(),
         TotalCount = totalCount,
         Page = request.Page,
         PageSize = request.PageSize
      };
   }

   public async Task<int> GetNotificationsCount(Guid? userId)
   {
      if (userId == null)
         return 0;

      return await _notificationRepository.CountNotifications(userId.Value,NotificationStatus.Unread);
   }

   public async Task DeleteNotification(Guid requestUserId, Guid notificationId)
   {
      var notification = await _notificationRepository.FindById(notificationId)
                         ?? throw new ApiException("Delete notification error", 404, "Notification wasn't found");
      
      if (notification.UserId != requestUserId)
      {
         throw new ApiException("Delete notification error", 403, "You can't delete this notification");
      }
      
      await _notificationRepository.DeleteAsync(notificationId);
      await UpdateNotificationsCount(requestUserId);
   }

   public async Task ChangeReadStatus(Guid userId, List<Guid> notificationIds)
   {
      await _notificationRepository.ChangeReadStatus(notificationIds);

      await UpdateNotificationsCount(userId);
   }
      

   public async Task<GetNotificationsCountResponse> GetUsersNotificationsCount(Guid userId)
   {
      return new GetNotificationsCountResponse
      {
         CountOfUnreadMessages = await _notificationRepository.CountNotifications(userId, NotificationStatus.Unread),
         CountOfReadMessages = await _notificationRepository.CountNotifications(userId, NotificationStatus.Read),
      };
   }

   public async Task<PagedResponse<Notification>> GetUserNotifications(Guid userId, PaginationRequest request)
   {
      var skip = (request.Page - 1) * request.PageSize;

      var totalCount = await _notificationRepository.CountAsync(n => n.UserId == userId);

      var notifications = await _notificationRepository.GetByPredicate(
         predicate: n => n.UserId == userId,
         orderBy: q => q.OrderByDescending(n => n.CreatedAt),
         skip: skip,
         take: request.PageSize
      );

      return new PagedResponse<Notification>
      {
         Items =notifications.ToList(),
         TotalCount = totalCount,
         Page = request.Page,
         PageSize = request.PageSize
      };
   }

   public async Task<string> GenerateOneTimeTicket(Guid userId)
   {
      var existResult = await _userClient.CheckIfUserExistAsync(new CheckIfUserExistsRequest
      {
         UserId = userId.ToString()
      });

      if (existResult.Value == false)
      {
         throw new ApiException("Generate one time ticket error", 404, "User wasn't found");
      }
      
      var userMappingKey = $"user_ticket_map_{userId}";

      if (_memoryCache.TryGetValue(userMappingKey, out string oldTicket))
      {
         _memoryCache.Remove($"ws_ticket_{oldTicket}");
      }

      var ticket = Guid.NewGuid().ToString("N");
      var expiration = TimeSpan.FromSeconds(30);

      _memoryCache.Set($"ws_ticket_{ticket}", userId.ToString(), expiration);
      _memoryCache.Set(userMappingKey, ticket, expiration);
      return ticket;
   }

   private async Task UpdateNotificationsCount(Guid userId)
   {
      var unreadMessagesCount = await _notificationRepository.CountNotifications(userId,NotificationStatus.Unread);
      await _hubContext.Clients.User(userId.ToString())
         .SendAsync(WebSocketMethodNames.UpdateUnreadMessagesMethod, unreadMessagesCount);
   }
}
   
