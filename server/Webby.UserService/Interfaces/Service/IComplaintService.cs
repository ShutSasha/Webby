using Webby.UserService.Dtos.Complaint;

namespace Webby.UserService.Interfaces.Service;

public interface IComplaintService
{
   Task CreateUserComplaint(CreateUserComplaintRequest request);
   Task<List<ComplaintDto>> GetUserComplaints(Guid userId);
}