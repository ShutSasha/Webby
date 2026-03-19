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
   private readonly IMapper _mapper;
   public UserService(IUserRepository userRepository, IMapper mapper, IStorageService storageService, IAchievementService achievementService)
   {
      _userRepository = userRepository;
      _mapper = mapper;
      _storageService = storageService;
      _achievementService = achievementService;
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

   public async Task FollowUser(UserFollowRequest request)
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

      if (userFollowExist)
      {
         throw new ApiException("Follow user error", 400, "You've already follow to this user");
      }

      await _userRepository.AddUserFollowing(request.UserId, request.FollowerId);
   }

   public async Task UnfollowUser(UserFollowRequest request)
   {
      if (request.UserId == request.FollowerId)
      {
         throw new ApiException("Unfollow user error", 400, "Invalid request");
      }
      
      var userFollowExist = await _userRepository.HasUserFollow(request.UserId, request.FollowerId);

      if (!userFollowExist)
      {
         throw new ApiException("Unfollow user error", 400, "Follow isn't exist");
      }

      await _userRepository.DeleteUserFollowing(request.UserId, request.FollowerId);
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
      //TODO: Add check user existance
      var userFollows = await _userRepository
         .GetUserFollowers(userId);

      return userFollows
         .Select(uf => _mapper.Map<UserFollowersDto>(uf.FollowerUser))
         .ToList();
   }

   public async Task<List<UserFollowersDto>> GetUserFollows(Guid userId)
   {
      //TODO: Add check user existance
      var userFollows = await _userRepository
         .GetUserFollows(userId);

      return userFollows
         .Select(uf => _mapper.Map<UserFollowersDto>(uf.FollowedUser))
         .ToList();
   }

   public async Task<User?> GetById(Guid userId) 
      => await _userRepository.FindById(userId);
   
}