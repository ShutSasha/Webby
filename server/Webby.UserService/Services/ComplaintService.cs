using AutoMapper;
using Grpc.Core;
using Webby.UserService.Clients;
using Webby.UserService.Dtos.Complaint;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;
using Webby.UserService.Models.Enums;

namespace Webby.UserService.Services;

public class ComplaintService : IComplaintService
{
   private readonly IComplaintRepository _complaintRepository;
   private readonly IUserService _userService;
   private readonly IMapper _mapper;
   private readonly VideoGrpcService.VideoGrpcServiceClient _videoGrpcServiceClient;
   
   public ComplaintService(IComplaintRepository complaintRepository, IUserService userService, IMapper mapper, VideoGrpcService.VideoGrpcServiceClient videoGrpcServiceClient)
   {
      _complaintRepository = complaintRepository;
      _userService = userService;
      _mapper = mapper;
      _videoGrpcServiceClient = videoGrpcServiceClient;
   }
   
   public async Task CreateComplaint(Guid authorId, CreateComplaintRequest request)
    {
        _ = await _userService.GetById(authorId) 
            ?? throw new ApiException("Create complaint error", 404, "Author user wasn't found");

        if (request.TargetType == null)
            throw new ApiException("Create complaint error", 400, "Target type is required");

        Guid parsedTargetId;

        switch (request.TargetType.Value)
        {
            case ComplaintTargetType.User:
                if (!Guid.TryParse(request.TargetId, out parsedTargetId))
                    throw new ApiException("Create complaint error", 400, "Invalid target user id format");

                if (authorId == parsedTargetId)
                    throw new ApiException("Create complaint error", 400, "The user cannot leave a complaint to himself");

                _ = await _userService.GetById(parsedTargetId) 
                    ?? throw new ApiException("Create complaint error", 404, "Target user wasn't found");
                break;

            case ComplaintTargetType.Video:
                var cleanVideoId = request.TargetId.StartsWith("wb_") 
                    ? request.TargetId[3..] 
                    : request.TargetId;

                if (!Guid.TryParse(cleanVideoId, out parsedTargetId))
                    throw new ApiException("Create complaint error", 400, "Invalid target video id format");

                try
                {
                    var checkVideoExistResult = await _videoGrpcServiceClient.CheckVideoExistsAsync(
                        new CheckVideoExistRequest
                        {
                            RequestUserId = authorId.ToString(),
                            VideoId = request.TargetId
                        });

                    if (!checkVideoExistResult.Value)
                        throw new ApiException("Create complaint error", 404, "Video wasn't found");
                }
                catch (RpcException ex) when (ex.StatusCode is StatusCode.NotFound or StatusCode.InvalidArgument)
                {
                    throw new ApiException("Create complaint error", 400, ex.Status.Detail);
                }
                catch (RpcException)
                {
                    throw new ApiException("Create complaint error", 503, "Video service is unavailable");
                }
                break;

            default:
                throw new ApiException("Create complaint error", 400, "Unknown target type");
        }

        var complaint = new Complaint
        {
            ComplaintId = Guid.NewGuid(),
            AdditionalInfo = request.AdditionalInfo,
            AuthorId = authorId,
            CreatedAt = DateTime.UtcNow,
            ReasonType = request.ReasonType,
            TargetId = parsedTargetId,
            TargetType = request.TargetType.Value
        };

        await _complaintRepository.Add(complaint);
    }

   public async Task<List<ComplaintDto>> GetUserComplaints(Guid userId)
   {
      var userComplaints = (await _complaintRepository
            .GetByPredicate(c => c.TargetId == userId))
         .ToList();
      
      return _mapper.Map<List<ComplaintDto>>(userComplaints);
   }

   public async Task<Complaint?> FindById(Guid complaintId)
   => await _complaintRepository.FindById(complaintId); 
}