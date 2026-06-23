using Webby.UserService.Dtos.Complaint;
using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Service;

public interface IComplaintService
{
   Task CreateComplaint(Guid authorId, CreateComplaintRequest request);
   Task<List<ComplaintDto>> GetUserComplaints(Guid userId);
   Task<Complaint?> FindById(Guid complaintId);
   
}