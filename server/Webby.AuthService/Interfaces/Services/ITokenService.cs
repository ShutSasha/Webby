using Webby.AuthService.Models;

namespace Webby.AuthService.Interfaces.Services;

public interface ITokenService
{
   Task<string> GenerateToken(User user);
   Task<Guid> ExtractUserInfo(string accessToken);

}