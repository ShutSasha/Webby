using AutoMapper;
using Webby.UserService.Dtos;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Services;

public class UserService : IUserService
{
   private readonly IUserRepository _userRepository;
   private readonly IMapper _mapper;
   public UserService(IUserRepository userRepository, IMapper mapper)
   {
      _userRepository = userRepository;
      _mapper = mapper;
   }
   
   public async Task<UserDto> GetUserInformation(Guid userId)
   {
      var user = await _userRepository.FindById(userId);

      if (user == null)
      {
         throw new ApiException("Get user information error", 404, "User wasn't found");
      }

      return _mapper.Map<UserDto>(user);
   }
   
}