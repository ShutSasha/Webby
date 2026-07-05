using Microsoft.EntityFrameworkCore;
using Webby.AchievementService.Data;
using Webby.AchievementService.Interfaces.Repositories;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Repositories;

public class EventTypeRepository : GenericRepository<EventType>, IEventTypeRepository
{
   public EventTypeRepository(AppDbContext context) : base(context)
   {
   }

   public async Task<bool> HasEventType(string eventName)
   {
      return await _context.EventTypes
         .AnyAsync(et => et.EventTypeName.ToLower() == eventName.ToLower());
   }
   
}