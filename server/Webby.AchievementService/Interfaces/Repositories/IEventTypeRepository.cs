using Webby.AchievementService.Models;

namespace Webby.AchievementService.Interfaces.Repositories;

public interface IEventTypeRepository : IRepository<EventType>
{
   Task<bool> HasEventType(string eventName);
}