using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class AchievementRepository : GenericRepository<Achievement>, IAchievementRepository
{
   public AchievementRepository(AppDbContext context) : base(context)
   {
   }
   
}