using Webby.UserService.Dtos;

namespace Webby.UserService.Interfaces.Service;

public interface IUserService
{
   Task<UserDto> GetUserInformation(Guid userId);
}