using AutoMapper;
using Webby.AuthService.Dtos;
using Webby.AuthService.Models;

namespace Webby.AuthService.Helpers.Mapping;

public class MappingProfiles: Profile
{
   public MappingProfiles()
   {
      CreateMap<User, UserDto>();
   }
}