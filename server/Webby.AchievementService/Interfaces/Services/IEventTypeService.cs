using Webby.AchievementService.Models;

namespace Webby.AchievementService.Interfaces.Services;

public interface IEventTypeService
{
   Task<Guid> AddEventType(string eventName);
   Task<List<EventType>> GetEventTypes();
   Task<bool> HasEventType(string eventName);
}