using Webby.UserService.Models;

namespace Webby.UserService.Interfaces.Repository;

public interface IComplaintRepository : IRepository<Complaint>
{
   Task SetIsBanComplaintStatus(Guid targetId, bool banStatusFlag);
}