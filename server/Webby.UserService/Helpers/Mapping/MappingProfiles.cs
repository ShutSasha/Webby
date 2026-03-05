using AutoMapper;
using Webby.UserService.Dtos;
using Webby.UserService.Models;

namespace Webby.UserService.Helpers.Mapping;

public class MappingProfiles : Profile
{
   public MappingProfiles()
   {
      CreateMap<User, UserDto>();

   }
}