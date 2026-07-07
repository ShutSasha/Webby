using Microsoft.EntityFrameworkCore;
using Webby.AuthService.Data;
using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Interfaces.Repositories;
using Webby.AuthService.Models;

namespace Webby.AuthService.Repositories;

public class AuthRepository : GenericRepository<User>, IAuthRepository
{
   public AuthRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<bool> HasUserRoles(Guid userId,IEnumerable<string> tokenRoles)
   {
      var user = await _context.Users.FirstOrDefaultAsync(u => u.UserId == userId);

      return user != null && tokenRoles.Contains(user.Role.ToString());
   }
}