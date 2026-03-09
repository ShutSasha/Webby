using AutoMapper;
using Webby.UserService.Dtos;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Services;

public class UserService : IUserService
{
   private readonly IUserRepository _userRepository;
   private readonly IStorageService _storageService;
   private readonly IMapper _mapper;
   public UserService(IUserRepository userRepository, IMapper mapper, IStorageService storageService)
   {
      _userRepository = userRepository;
      _mapper = mapper;
      _storageService = storageService;
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
}