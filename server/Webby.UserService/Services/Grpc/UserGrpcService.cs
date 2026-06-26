using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using UserService;
using Webby.NotificationService.GrpcClient;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models.Enums;

namespace Webby.UserService.Services.Grpc;

public class UserGrpcService : global::UserService.UserGrpcService.UserGrpcServiceBase
{
   private readonly IUserService _userService;
   private readonly IUserRepository _userRepository;
   private readonly IUserPremiumRepository _userPremiumRepository;
   private readonly IPaymentRepository _paymentRepository;
   private readonly NotificationGrpcService.NotificationGrpcServiceClient _notificationClient;
   private readonly ILogger<UserGrpcService> _logger;

   public UserGrpcService(IUserRepository userRepository, ILogger<UserGrpcService> logger,
      IUserPremiumRepository userPremiumRepository, NotificationGrpcService.NotificationGrpcServiceClient notificationClient, IPaymentRepository paymentRepository, IUserService userService)
   {
      _userRepository = userRepository;
      _logger = logger;
      _userPremiumRepository = userPremiumRepository;
      _notificationClient = notificationClient;
      _paymentRepository = paymentRepository;
      _userService = userService;
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
      if (!Guid.TryParse(request.UserID, out var userIdGuid)
          || !Guid.TryParse(request.RequestUserID, out var requestUserIdGuid))
         throw new ApiException("Ban user error",400,"Invalid id format");

      await _userService.BanUser(requestUserIdGuid, userIdGuid);

      return new Empty();
   }

   public override async Task<GetMonthlyRegistrationsResponse> GetMonthlyRegistrations(Empty request, ServerCallContext context)
   {
      var (startDate, endDate) = GetYearDates();
      
      var monthlyData = new Dictionary<int, int>();
      for (var i = 0; i < 12; i++)
      {
         var targetMonth = startDate.AddMonths(i).Month;
         monthlyData[targetMonth] = 0;
      }
      
      var registrationsFromDb = await _userRepository.GetMonthlyRegistrationsCountAsync(startDate, endDate);
      
      foreach (var reg in registrationsFromDb)
      {
         monthlyData[reg.Key] = reg.Value;
      }

      var response = new GetMonthlyRegistrationsResponse();
      
      foreach (var kvp in monthlyData)
      {
         response.Registrations.Add(kvp.Key, kvp.Value);
      }

      return response;
   }

   public override async Task<GetMonthlySubscriptionsResponse> GetMonthlySubscriptions(Empty request, ServerCallContext context)
   {
      var (startDate, endDate) = GetYearDates();
      var response = new GetMonthlySubscriptionsResponse();
      
      for (var i = 0; i < 12; i++)
      {
         var targetMonth = startDate.AddMonths(i).Month;
         response.Subscriptions[targetMonth] = 0;
      }
      
      var revenueInDb = await _paymentRepository.GetMonthlySubscriptionRevenueAsync(startDate, endDate);
      
      foreach (var kvp in revenueInDb)
      {
         response.Subscriptions[kvp.Key] = (int)kvp.Value; 
      }

      return response;
   }

   public override async Task<GetTotalRegistrationsResponse> GetTotalRegistrations(Empty request, ServerCallContext context)
   {
      var totalRegistrationsCount = await _userRepository.CountAsync();
      return new GetTotalRegistrationsResponse() { TotalRegistrations = totalRegistrationsCount };
   }

   public override async Task<GetMonthRevenueResponse> GetMonthRevenue(Empty request, ServerCallContext context)
   {
      var monthRevenue = await _paymentRepository.GetMonthRevenue();
      return new GetMonthRevenueResponse() { MonthRevenue = monthRevenue };
   }

   public override async Task<BoolValue> CheckUserCanBlockVideo(CheckUserCanBlockVideoRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.VideoAuthorId, out var videoAuthorIdGuid)
          || !Guid.TryParse(request.RequestedUserId, out var requestedUserIdGuid))
      {
         throw new ApiException("Check user can block error", 400, "Invalid id format type");
      }

      var canChangeVideoBanStatus = await CanChangeVideoBanStatus(requestedUserIdGuid, videoAuthorIdGuid);

      return new BoolValue { Value = canChangeVideoBanStatus };

   }

   
   private async Task<bool> CanChangeVideoBanStatus(Guid requestedUserId,Guid authorId)
   {
      var requester = await _userRepository.FindById(requestedUserId);
      var author = await _userRepository.FindById(authorId);

      if (requester == null || author == null)
      {
         throw new ApiException("Check user can block error", 404, "User wasn't found");
      }
      
      if (requestedUserId == authorId) 
      {
         return false; 
      }

      return requester.Role switch
      {
         Role.Admin => author.Role != Role.Admin,
         Role.Moderator => author.Role == Role.User, 
         _ => false 
      };
   }
   
   private Tuple<DateTime,DateTime > GetYearDates()
   {
      var today = DateTime.UtcNow;

      var startDate = new DateTime(today.Year, today.Month, 1, 0, 0, 0, DateTimeKind.Utc).AddMonths(-11);
      var endDate = new DateTime(today.Year, today.Month, 1, 0, 0, 0, DateTimeKind.Utc).AddMonths(1);

      return new Tuple<DateTime, DateTime>(startDate, endDate);
   }
}