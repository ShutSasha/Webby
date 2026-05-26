using AutoMapper;
using Webby.AchievementService.Dtos.Achievement;
using Webby.AchievementService.Models;

namespace Webby.AchievementService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<Achievement, AchievementDto>();
      CreateMap<Achievement, ProfileAchievementDto>();
   }
}