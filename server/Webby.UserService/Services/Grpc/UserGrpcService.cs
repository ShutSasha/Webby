using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using UserService;
using Webby.NotificationService.GrpcClient;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;

namespace Webby.UserService.Services.Grpc;

public class UserGrpcService : global::UserService.UserGrpcService.UserGrpcServiceBase
{
   private readonly IUserRepository _userRepository;
   private readonly IUserPremiumRepository _userPremiumRepository;
   private readonly NotificationGrpcService.NotificationGrpcServiceClient _notificationClient;
   private readonly ILogger<UserGrpcService> _logger;

   public UserGrpcService(IUserRepository userRepository, ILogger<UserGrpcService> logger,
      IUserPremiumRepository userPremiumRepository, NotificationGrpcService.NotificationGrpcServiceClient notificationClient)
   {
      _userRepository = userRepository;
      _logger = logger;
      _userPremiumRepository = userPremiumRepository;
      _notificationClient = notificationClient;
   }

   public override async Task<UserResponse> GetUserById(GetUserRequest request, ServerCallContext context)
   {
      var user = await _userRepository.FindById(Guid.Parse(request.UserId));
      
      if (user == null)
         throw new RpcException(new Status(StatusCode.NotFound, "User not found"));
      
      var response = new UserResponse
      {
         UserId = user.UserId.ToString(),
         Username = user.Username,
         AvatarUrl = user.AvatarUrl,
         IsFollowed = request.RequestUserId != string.Empty && await _userRepository
            .HasUserFollow(Guid.Parse(request.UserId),Guid.Parse(request.RequestUserId))
      };
      
      return response;
   }

   public override async Task<GetUsersResponse> GetUsersByIds(GetUsersRequest request, ServerCallContext context)
   {
      var userIds = request.UserIds
         .Select(Guid.Parse)
         .ToList();

      var orderMap = userIds
         .Select((id, index) => (id, index))
         .ToDictionary(x => x.id, x => x.index);

      var users = await _userRepository.GetByIds(userIds);

      var mappedUsers = users
         .OrderBy(u => orderMap[u.UserId])
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

   public override async Task<GetUserSubscriptionsIdsResponse> GetUserSubscriptionIds(GetUserSubscriptionIdsRequest request, ServerCallContext context)
   {
      var userSubscriptionIds = await _userRepository.GetSubscriptionIds(Guid.Parse(request.RequestUserId));

      return new GetUserSubscriptionsIdsResponse()
      {
         UserIds = { userSubscriptionIds.Select(u => u.ToString()) }
      };
   }

   public override async Task<GetPremiumStatusResponse> GetPremiumStatus(GetPremiumStatusRequest request, ServerCallContext context)
   {
      var userPremium = await _userPremiumRepository.GetUserPremiumInformation(Guid.Parse(request.UserId));
   
      var status = userPremium switch
      {
         null => PremiumStatus.None,
         _ when userPremium.ExpiresAt > DateTime.UtcNow => PremiumStatus.Active,
         _ => PremiumStatus.Expired
      };

      return new GetPremiumStatusResponse
      {
         Status = status
      };
   }

   public override async Task<BoolValue> CheckIfUserExist(CheckIfUserExistsRequest request, ServerCallContext context)
   {
      var user = await _userRepository.FindById(Guid.Parse(request.UserId));

      return user == null ? new BoolValue { Value = false } : new BoolValue { Value = true };
   }
   
   public override async Task<FindUserIDsResponse> FindUserIDs(FindUserIDsRequest request, ServerCallContext context)
   {
      var userIdGuids = request.UserIds.Select(Guid.Parse).ToList();
      var userResult = await _userRepository.SearchByUsername(userIdGuids, request.Search, request.Limit, request.Offset);

      return new FindUserIDsResponse
      {
         UserIds = {userResult.Ids.Select(id => id.ToString())},
         Total = userResult.Total
      };
      
   }

   public override async Task<IsFollowedResponse> IsFollowed(IsFollowedRequest request, ServerCallContext context)
   {
      var mutualFollows =
         await _userRepository.AreMutualFollowers(Guid.Parse(request.FirstUserId), Guid.Parse(request.SecondUserId));

      return new IsFollowedResponse
      {
         IsFollowed = mutualFollows
      };
   }

   public override async Task<Empty> BanUser(BanUserRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.UserID, out var userIdGuid))
         throw new ApiException("Ban user error",400,"Invalid id format");

      var user = await _userRepository.FindById(userIdGuid)
                 ?? throw new ApiException("Ban user error", 404,"User wasn't found");


      if (user.IsBanned)
         return new Empty();

      user.IsBanned = true;

      await _userRepository.Update(user);

      await _notificationClient.ReportBlockingAsync(new ReportBlockingRequest { UserId = request.UserID });

      return new Empty();
   }
   
}