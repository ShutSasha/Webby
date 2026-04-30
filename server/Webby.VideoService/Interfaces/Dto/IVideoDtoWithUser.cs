using Webby.VideoService.Dtos.User;

namespace Webby.VideoService.Interfaces.Dto;

public interface IVideoDtoWithUser
{
   UserVideoDto? User { get; set; }
}