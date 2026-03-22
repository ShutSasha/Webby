using Webby.UserService.Dtos.Complaint;

namespace Webby.UserService.Interfaces.Service;

public interface IComplaintService
{
   Task CreateComplaint(Guid authorId, CreateComplaintRequest request);
   Task<List<ComplaintDto>> GetUserComplaints(Guid userId);
}