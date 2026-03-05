using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class ComplaintRepository : GenericRepository<Complaint>, IComplaintRepository
{
   public ComplaintRepository(AppDbContext context) : base(context)
   {
   }
}