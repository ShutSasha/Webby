using Grpc.Core;
using UserService;
using Webby.UserService.Interfaces.Repository;

namespace Webby.UserService.Services;

public class UserGrpcService : global::UserService.UserGrpcService.UserGrpcServiceBase
{
   private readonly IUserRepository _userRepository;
   private readonly ILogger<UserGrpcService> _logger;

   public UserGrpcService(IUserRepository userRepository, ILogger<UserGrpcService> logger)
   {
      _userRepository = userRepository;
      _logger = logger;
   }

   // public override async Task<UserResponse> GetUserById(GetUserRequest request, ServerCallContext context)
   // {
   //    // var user = await _userRepository.FindById(Guid.Parse(request.UserId));
   //    // if (user == null)
   //    //    throw new RpcException(new Status(StatusCode.NotFound, "User not found"));
   //    //
   //    // var response = new UserResponse
   //    // {
   //    //    UserId = user.UserId.ToString(),
   //    //    Username = user.Username,
   //    //    AvatarUrl = user.AvatarUrl
   //    // };
   //    // Console.WriteLine($"RESPONSE FROM GRPC SERVER {response} ");
   //    //
   //    // return response;
   //    return new UserResponse { UserId = "1", Username = "Test", AvatarUrl = "" };
   // }
   
   public override async Task<UserResponse> GetUserById(GetUserRequest request, ServerCallContext context)
   {
      var response = new UserResponse { UserId = "1", Username = "Test", AvatarUrl = "" };
      return await Task.FromResult(response);
   }

   public override async Task<GetUsersResponse> GetUsersByIds(GetUsersRequest request, ServerCallContext context)
   {
      var userIds = request.UserIds
         .Select(Guid.Parse)
         .ToList();

      var users = await _userRepository.GetByIds(userIds);
      
      var mappedUsers = users
         .Select(u => new UserResponse
         {
            UserId = u.UserId.ToString(),
            Username = u.Username,
            AvatarUrl = u.AvatarUrl
         })
         .ToList();

      return new GetUsersResponse
      {
         Users = { mappedUsers }
      };
      
   }
}