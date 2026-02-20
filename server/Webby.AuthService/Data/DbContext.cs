using Microsoft.EntityFrameworkCore;
using Webby.AuthService.Models;

namespace Webby.AuthService.Data;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
   public DbSet<User> Users { get; set; }
}