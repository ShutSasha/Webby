using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Webby.UserService.ComplaintGrpcService;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Service;

namespace Webby.UserService.Services.Grpc;

public class ComplaintGrpcService : Webby.UserService.ComplaintGrpcService.ComplaintGrpcService.ComplaintGrpcServiceBase
{
   private readonly IComplaintService _complaintService;

   public ComplaintGrpcService(IComplaintService complaintService)
   {
      _complaintService = complaintService;
   }
   
   public override async Task<GetComplaintByIDResponse> GetComplaintByID(GetComplaintByIDRequest request, ServerCallContext context)
   {
      var complaint = await _complaintService.FindById(Guid.Parse(request.ComplaintID));

      if (complaint == null)
      {
         throw new ApiException("Get complaint error", 404, "Complaint wasn't found");
      }

      return new GetComplaintByIDResponse
      {
         ComplaintID = complaint.ComplaintId.ToString(),
         AuthorId = complaint.AuthorId.ToString(),
         TargetType = complaint.TargetType.ToString(),
         TargetId = complaint.TargetId.ToString(),
         ReasonType = complaint.ReasonType,
         AdditionalInfo = complaint.AdditionalInfo,
         CreatedAt = complaint.CreatedAt.ToTimestamp()
      };
   }

   public override async Task<Empty> SetVideoComplaintsBanned(SetVideoComplaintsBannedRequest request, ServerCallContext context)
   {
      if (!Guid.TryParse(request.TargetId, out var targetIdGuid))
      {
         throw new ApiException("Send video complaints banned", 400, "Incorrect type of target id");
      }

      await _complaintService.SetBanStatus(targetIdGuid, request.BanFlag);

      return new Empty();
   }
}