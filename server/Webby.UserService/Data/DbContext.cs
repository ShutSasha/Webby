using Microsoft.EntityFrameworkCore;
using Webby.UserService.Models;

namespace Webby.UserService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<User> Users { get; set; }
}