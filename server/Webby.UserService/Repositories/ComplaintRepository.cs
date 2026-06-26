using Microsoft.EntityFrameworkCore;
using Webby.UserService.Data;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Models;

namespace Webby.UserService.Repositories;

public class ComplaintRepository : GenericRepository<Complaint>, IComplaintRepository
{
   public ComplaintRepository(AppDbContext context) : base(context)
   {
   }

   public async Task SetIsBanComplaintStatus(Guid targetId,bool banStatusFlag)
   {
      await _context.Complaints
         .Where(c => c.ComplaintId == targetId)
         .ExecuteUpdateAsync(c => c.SetProperty(c => c.IsBanned, banStatusFlag));
   }
}