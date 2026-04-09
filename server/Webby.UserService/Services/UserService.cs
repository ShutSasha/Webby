using AutoMapper;
using Webby.UserService.Dtos;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;

namespace Webby.UserService.Services;

public class UserService : IUserService
{
   private readonly IUserRepository _userRepository;
   private readonly IStorageService _storageService;
   private readonly IAchievementService _achievementService;
   private readonly IUserPremiumRepository _userPremiumRepository;
   private readonly IMapper _mapper;
   public UserService(IUserRepository userRepository, IMapper mapper, IStorageService storageService,
      IAchievementService achievementService, IUserPremiumRepository userPremiumRepository)
   {
      _userRepository = userRepository;
      _mapper = mapper;
      _storageService = storageService;
      _achievementService = achievementService;
      _userPremiumRepository = userPremiumRepository;
   }
   
   public async Task<UserProfileResponse> GetUserInformation(Guid userId)
   {
      var response = new UserProfileResponse();
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Get user information error", 404, "User wasn't found");
      }
      var followBlock = await _userRepository.GetUserFollowBlock(userId);

      response.User = _mapper.Map<UserDto>(user);
      response.UserFollowStats = followBlock;
      response.PinnedUserAchievements = await _achievementService.GetPinnedAchievements(userId);
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

      if (user.AvatarUrl.StartsWith("https://webby-watch-platform-bucket"))
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
            return "Successfully followed user";
      }
      
   }

   public async Task UnlockAchievement(Guid userId, Guid achievementId)
   {
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Unlock achievement error", 404, "User wasn't found");
      }

      var achievement = await _achievementService.FindById(achievementId);

      if (achievement == null)
      {
         throw new ApiException("Unlock achievement error", 404, "Achievement wasn't found");
      }

      await _achievementService.AddUserAchievement(userId, achievementId);
      
   }

   public async Task PinUserAchievement(Guid userId, Guid achievementId)
   {
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Pin achievement error", 404, "User wasn't found");
      }
      
      var achievement = await _achievementService.FindById(achievementId);

      if (achievement == null)
      {
         throw new ApiException("Pin achievement error", 404, "Achievement wasn't found");
      }
      
      var userAchievement = await _achievementService.GetUserAchievement(userId,achievementId);

      if (await _achievementService.GetPinnedAchievementsCount(user.UserId) >= 3)
      {
         throw new ApiException("Pin user achievement error", 400, "You can't pin more than 3 achievements");
      }

      userAchievement.IsPinned = true;
      await _achievementService.UpdateUserAchievement(userAchievement);
   }

   public async Task UnpinUserAchievement(Guid userId, Guid achievementId)
   {
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Unpin achievement error", 404, "User wasn't found");
      }
      
      var achievement = await _achievementService.FindById(achievementId);

      if (achievement == null)
      {
         throw new ApiException("Unpin achievement error", 404, "Achievement wasn't found");
      }

      var userAchievement = await _achievementService.GetUserAchievement(userId,achievementId);

      userAchievement.IsPinned = false;
      await _achievementService.UpdateUserAchievement(userAchievement);
      

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
}