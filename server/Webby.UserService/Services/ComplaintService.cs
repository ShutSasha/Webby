using AutoMapper;
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
   
   public ComplaintService(IComplaintRepository complaintRepository, IUserService userService, IMapper mapper)
   {
      _complaintRepository = complaintRepository;
      _userService = userService;
      _mapper = mapper;
   }
   
   public async Task CreateUserComplaint(CreateUserComplaintRequest request)
   {
      if (request.AuthorId == request.TargetUserId)
         throw new ApiException("Create complaint error", 400, "The user cannot leave a complaint to himself");
      
      var user = await _userService.GetById(request.AuthorId);

      if (user == null)
      {
         throw new ApiException("Create complaint error", 404, "User wasn't found");
      }

      var targetUser = await _userService.GetById(request.TargetUserId);

      if (targetUser == null)
      {
         throw new ApiException("Create complaint error", 404, "Target user wasn't found");
      }

      var userComplaint = new Complaint()
      {
         ComplaintId = Guid.NewGuid(),
         AdditionalInfo = request.AdditionalInfo,
         AuthorId = request.AuthorId,
         CreatedAt = DateTime.UtcNow,
         ReasonType = request.ReasonType,
         TargetId = request.TargetUserId,
         TargetType = ComplaintTargetType.User
      };

      await _complaintRepository.Add(userComplaint);
   }

   public async Task<List<ComplaintDto>> GetUserComplaints(Guid userId)
   {
      var userComplaints = (await _complaintRepository
            .GetByPredicate(c => c.TargetId == userId))
         .ToList();
      
      return _mapper.Map<List<ComplaintDto>>(userComplaints);
   }
   
}