using AutoMapper;
using Webby.UserService.Dtos;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Dtos.Complaint;
using Webby.UserService.Dtos.User;
using Webby.UserService.Models;

namespace Webby.UserService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<User, UserDto>();
      CreateMap<User, UserFollowersDto>();
      CreateMap<Complaint, ComplaintDto>();
   }
}