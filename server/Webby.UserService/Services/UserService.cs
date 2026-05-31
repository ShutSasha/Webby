using AutoMapper;
using Grpc.Core;
using UserService.AchievementGrpcClient;
using Webby.NotificationService.GrpcClient;
using Webby.UserService.Consts;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Dtos.User;
using Webby.UserService.Dtos.Notification;
using Webby.UserService.Dtos.Search;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Helpers.Response;
using Webby.UserService.Interfaces.Helpers;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;

namespace Webby.UserService.Services;

public class UserService : IUserService
{
   private readonly IUserRepository _userRepository;
   private readonly IStorageService _storageService;
   private readonly IUserPremiumRepository _userPremiumRepository;
   private readonly NotificationGrpcService.NotificationGrpcServiceClient _notificationGrpcServiceClient;
   private readonly AchievementGrpcService.AchievementGrpcServiceClient _achievementGrpcServiceClient;
   private readonly INotificationFactory _notificationFactory;
   private readonly ILogger<UserService> _logger;
   
   
   private readonly IMapper _mapper;
   public UserService(IUserRepository userRepository, IMapper mapper, IStorageService storageService,
      IUserPremiumRepository userPremiumRepository, NotificationGrpcService.NotificationGrpcServiceClient notificationGrpcServiceClient,
      INotificationFactory notificationFactory,AchievementGrpcService.AchievementGrpcServiceClient achievementGrpcServiceClient,
      ILogger<UserService> logger)
   {
      _userRepository = userRepository;
      _mapper = mapper;
      _storageService = storageService;
      _userPremiumRepository = userPremiumRepository;
      _notificationGrpcServiceClient = notificationGrpcServiceClient;
      _notificationFactory = notificationFactory;
      _achievementGrpcServiceClient = achievementGrpcServiceClient;
      _logger = logger;
   }
   
   public async Task<UserProfileResponse> GetUserInformation(Guid userId)
   {
      var response = new UserProfileResponse();
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Get user information error", 404, "User wasn't found");
      }
      
      var followBlockTask = _userRepository.GetUserFollowBlock(userId);
      

      response.User = _mapper.Map<UserDto>(user);

      try
      {
         var userAchievements = await _achievementGrpcServiceClient
            .GetPinnedAchievementsAsync(new GetPinnedAchievementsRequest { UserId = userId.ToString() });

         response.PinnedUserAchievements = userAchievements.UserAchievements
            .Select(ua => new ProfileAchievementDto
            {
               AchievementId = Guid.Parse(ua.AchievementId),
               IconUrl = ua.IconUrl,
               Title = ua.Title

            }).ToList();
      }
      catch (RpcException ex)
      {
         _logger.LogWarning(ex, "Failed to fetch pinned achievements via gRPC for user {UserId}. Status: {StatusCode}", userId, ex.StatusCode);
         response.PinnedUserAchievements = [];
      }

      response.UserFollowStats = await followBlockTask;
      
      return response;
   }

   public async Task<UserDto> UpdateUserInformation(UpdateUserRequest request)
   {
      var user = await _userRepository.FindById(request.UserId);
      
      if (user == null)
      {
         throw new ApiException("Update user information error", 404, "User wasn't found");
      }

      user.About = request.About ?? string.Empty;

      await _userRepository.Update(user);
      
      return _mapper.Map<UserDto>(user);
   }

   public async Task<UserDto> EditUserIcon(Guid id, string fileName, Stream fileStream, string contentType)
   {
      var user = await _userRepository.FindById(id);
      
      if (user == null)
      {
         throw new ApiException("Update user information error", 404, "User wasn't found");
      }

      var updatedUserIconPath = await _storageService
         .UploadFileAsync(id,"user_data",fileName,fileStream,contentType);

      if (user.AvatarUrl.StartsWith("https://webby-watch-platform-bucket") && user.AvatarUrl != DefaultLinks.DefaultUserIcon)
      {
         await _storageService.DeleteFileAsync(user.AvatarUrl);
      }

      user.AvatarUrl = updatedUserIconPath;

      await _userRepository.Update(user);

      return _mapper.Map<UserDto>(user);

   }

   public async Task<string> ProcessFollow(UserFollowRequest request)
   {
      if (request.UserId == request.FollowerId)
      {
         throw new ApiException("Follow user error", 400, "User cannot follow himself");
      }

      var user = await _userRepository.FindById(request.UserId);

      if (user == null)
      {
         throw new ApiException("Follow user error", 404, "User wasn't found");
      }

      var userFollowExist = await _userRepository.HasUserFollow(request.UserId, request.FollowerId);

      switch (userFollowExist)
      {
         case true:
            await _userRepository.DeleteUserFollowing(request.UserId, request.FollowerId);
            return "Successfully unfollowed user";

         case false:
            await _userRepository.AddUserFollowing(request.UserId, request.FollowerId);
            await SendNotification(await _notificationFactory.CreateNewFollowerNotification(request.UserId,request.FollowerId));
            return "Successfully followed user";
      }
      
   }

   public async Task PinUserAchievement(Guid userId, Guid achievementId)
   {
      var user = await _userRepository.FindById(userId);
   
      if (user == null)
      {
         throw new ApiException("Pin achievement error", 404, "User wasn't found");
      }
      
      var request = new ProcessPinAchievementsRequest
      {
         UserId = userId.ToString(),
         AchievementId = achievementId.ToString(),
         ProcessAchievementType = ProcessAchievementType.Pin
      };

      await _achievementGrpcServiceClient.ProcessPinAchievementAsync(request);
   }
   
   public async Task UnpinUserAchievement(Guid userId, Guid achievementId)
   {
      var user = await _userRepository.FindById(userId);
   
      if (user == null)
      {
         throw new ApiException("Unpin achievement error", 404, "User wasn't found");
      }
      
      var request = new ProcessPinAchievementsRequest
      {
         UserId = userId.ToString(),
         AchievementId = achievementId.ToString(),
         ProcessAchievementType = ProcessAchievementType.Unpin
      };

      await _achievementGrpcServiceClient.ProcessPinAchievementAsync(request);
   }

   public async Task<List<UserFollowersDto>> GetUserFollowers(Guid userId)
   {
      _ = await _userRepository.FindById(userId)
          ?? throw new ApiException("Get user followers error", 404, "User wasn't found");
      
      var userFollowers = await _userRepository
         .GetUserFollowers(userId);

      return userFollowers
         .Select(uf => _mapper.Map<UserFollowersDto>(uf.FollowerUser))
         .ToList();
   }

   public async Task<List<UserFollowersDto>> GetUserFollows(Guid userId)
   {
      _ = await _userRepository.FindById(userId)
          ?? throw new ApiException("Get user followers error", 404, "User wasn't found");
      
      var userFollows = await _userRepository
         .GetUserFollows(userId);

      return userFollows
         .Select(uf => _mapper.Map<UserFollowersDto>(uf.FollowedUser))
         .ToList();
   }

   public async Task<User?> GetById(Guid userId) 
      => await _userRepository.FindById(userId);

   public async Task<UserFollowingResponse> IsUserFollowing(Guid userId, Guid targetId)
   {
      var response = new UserFollowingResponse();
      
      var targetUser = await _userRepository.FindById(targetId);

      if (targetUser == null)
         throw new ApiException("Check is user following error", 404, "Target user wasn't found");

      response.isFollowing = await _userRepository.HasUserFollow(targetId, userId);
      return response;
   }

   public async Task<GetUserSubscriptionResponse?> GetUserSubscription(Guid userId, bool isOwner)
   {
      if (!isOwner)
      {
         throw new ApiException("Get user subscription error", 403,
            "You don't have permission to see subscription of this user");
      }

      var userPremiumInformation = await _userPremiumRepository.GetUserPremiumInformation(userId);

      if (userPremiumInformation == null)
      {
         return null;
      }

      return new GetUserSubscriptionResponse()
      {
         ExpiresAt = userPremiumInformation.ExpiresAt
      };
      
   }

   public async Task<PagedResponse<UserDto>> SearchUsers(SearchOptions searchOptions)
   {
      var skip = (searchOptions.Page - 1) * searchOptions.PageSize;

      var (users, total) = await _userRepository.SearchAsync(
         "Users",
         "Username",
         searchOptions.SearchText,
         skip,
         searchOptions.PageSize);

      return new PagedResponse<UserDto>()
      {
         Items = users.Count != 0 ? users.Select(u => _mapper.Map<UserDto>(u)).ToList() : [],
         TotalCount = total,
         Page = searchOptions.Page,
         PageSize = searchOptions.PageSize
      };
   }

   public async Task<Guid> DeleteUser(Guid userId)
   {
      var user = await _userRepository.FindById(userId) ??
                 throw new ApiException("Delete user error", 404, "User wasn't found");

      return await _userRepository.DeleteAsync(userId);
   }

   private async Task SendNotification(SendNotificationDto notificationDto) 
      => await _notificationGrpcServiceClient.SendNotificationToUserAsync(new CreateNotificationRequest
      {
         Message = notificationDto.Message,
         TargetType = GrpcNotificationTargetType.User,
         TargetIdentifier = notificationDto.TargetIdentifier,
         Title = notificationDto.Title,
         UserId = notificationDto.UserId
      });
}