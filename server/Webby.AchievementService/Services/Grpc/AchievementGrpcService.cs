using AchievementService;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Webby.AchievementService.Helpers.Exception;
using Webby.AchievementService.Interfaces.Services;

namespace Webby.AchievementService.Services.Grpc;

public class AchievementGrpcService : global::AchievementService.AchievementGrpcService.AchievementGrpcServiceBase
{
   private readonly IAchievementService _achievementService;

   public AchievementGrpcService(IAchievementService achievementService)
   {
      _achievementService = achievementService;
   }
   
   public override async Task<BoolValue> ProcessPinAchievement(ProcessPinAchievementsRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.AchievementId, out var achievementId) || 
          !Guid.TryParse(request.UserId, out var userId))
      {
         throw new ApiException("Process achievement error", 400, "Invalid format for UserId or AchievementId. Must be a UUID.");
      }

      var achievement = await _achievementService.FindById(achievementId);
      if (achievement == null)
      {
         throw new ApiException("Process achievement error", 404, "Achievement wasn't found");
      }

      var userAchievement = await _achievementService.GetUserAchievement(userId, achievementId);
      if (userAchievement == null)
      {
         throw new ApiException("Process achievement error", 400, "User has not unlocked this achievement yet");
      }

      switch (request.ProcessAchievementType)
      {
         case ProcessAchievementType.Unspecified:
            throw new ApiException("Process achievement error", 400, "Unspecified type of process achievement");
      
         case ProcessAchievementType.Pin:
            if (userAchievement.IsPinned) 
               break;

            if (await _achievementService.GetPinnedAchievementsCount(userId) >= 3)
            {
               throw new ApiException("Pin user achievement error", 400, "You can't pin more than 3 achievements");
            }

            userAchievement.IsPinned = true;
            await _achievementService.UpdateUserAchievement(userAchievement);
            break;

         case ProcessAchievementType.Unpin:
            if (!userAchievement.IsPinned) 
               break;

            userAchievement.IsPinned = false;
            await _achievementService.UpdateUserAchievement(userAchievement);
            break;
      }

      return new BoolValue { Value = true };
   }

   public override async Task<GetPinnedAchievementsResponse> GetPinnedAchievements(GetPinnedAchievementsRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.UserId, out var userId))
      {
         throw new ApiException("Get pinned achievements error", 400, "Incorrect format of user id");
      }

      var userAchievements = await _achievementService.GetPinnedAchievements(userId);
      var response = new GetPinnedAchievementsResponse();
      var achievementItems = userAchievements.Select(ua => new ProfileAchievementItem
      {
         AchievementId = ua.AchievementId.ToString(),
         IconUrl = ua.IconUrl,
         Title = ua.Title

      });
      
      response.UserAchievements.AddRange(achievementItems);
      return response;
   }
   
}