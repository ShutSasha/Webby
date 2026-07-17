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
      CreateMap<UserAchievement,AchievementDto>()
         .ForMember(dest => dest.AchievementId, opt => opt.MapFrom(src => src.AchievementId))
         .ForMember(dest => dest.Title, opt => opt.MapFrom(src => src.Achievement.Title))
         .ForMember(dest => dest.Description, opt => opt.MapFrom(src => src.Achievement.Description))
         .ForMember(dest => dest.IconUrl, opt => opt.MapFrom(src => src.Achievement.IconUrl))
         .ForMember(dest => dest.TargetValue, opt => opt.MapFrom(src => src.Achievement.TargetValue))
         .ForMember(dest => dest.IsUnlocked, opt => opt.MapFrom(src => true))
         .ForMember(dest => dest.UnlockedAt, opt => opt.MapFrom(src => src.UnlockedAt))
         .ForMember(dest => dest.AchievementProgressValue, opt => opt.MapFrom(src => src.Achievement.TargetValue));
   }
}