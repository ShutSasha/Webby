using Webby.AchievementService.Interfaces.Repositories;
using Webby.AchievementService.Interfaces.Services;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Services;

public class EventTypeService : IEventTypeService
{
   private readonly IEventTypeRepository _eventTypeRepository;
   public EventTypeService(IEventTypeRepository eventTypeRepository)
   {
      _eventTypeRepository = eventTypeRepository;
   }
   
   public async Task<Guid> AddEventType(string eventName)
   {
      var eventType = new EventType()
      {
         EventTypeId = Guid.NewGuid(),
         EventTypeName = eventName
      };

      var eventTypeId = await _eventTypeRepository.Add(eventType);
      return eventTypeId;
   }

   public async Task<List<EventType>> GetEventTypes()
      => await _eventTypeRepository.GetAll();

   public async Task<bool> HasEventType(string eventName)
      => await _eventTypeRepository.HasEventType(eventName);
}