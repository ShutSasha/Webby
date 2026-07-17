using System.Linq.Expressions;
using Microsoft.EntityFrameworkCore;
using Webby.NotificationService.Data;
using Webby.NotificationService.Interfaces.Repositories;

namespace Webby.NotificationService.Repositories;

public class GenericRepository<T> : IRepository<T> where T : class
{
   protected readonly  AppDbContext _context;
   private readonly DbSet<T> _dbSet;

   public GenericRepository(AppDbContext context)
   {
      _context = context;
      _dbSet = _context.Set<T>();
   }
   
   public async Task<List<T>> GetAll()
   {
      return await _dbSet.ToListAsync();
   }
   
   public async Task<Guid> Add(T entity)
   {
      if (entity == null)
      {
         throw new ArgumentNullException(nameof(entity));
      }

      await _dbSet.AddAsync(entity);
      await _context.SaveChangesAsync();
      
      var property = entity.GetType().GetProperty("Id") ?? entity.GetType().GetProperty("EntityId");
      if (property != null && property.PropertyType == typeof(Guid))
      {
         return (Guid)property.GetValue(entity);
      }

      return Guid.Empty;
   }
   
   public async Task<T?> FindById(Guid id)
   {
      return await _dbSet.FindAsync(id);
   }
   
   public async Task Update(T entity)
   {
      if (entity == null)
      {
         throw new ArgumentNullException(nameof(entity));
      }

      _dbSet.Update(entity);
      await _context.SaveChangesAsync();
   }
   
   public async Task<Guid> DeleteAsync(Guid id)
   {
      var entity = await FindById(id);
      if (entity == null)
      {
         throw new KeyNotFoundException($"Entity with ID {id} not found.");
      }

      _dbSet.Remove(entity);
      await _context.SaveChangesAsync();
      return id;
   }
   public async Task<IEnumerable<T>> GetByPredicate(Expression<Func<T, bool>> predicate, 
      Func<IQueryable<T>, IOrderedQueryable<T>>? orderBy = null,
      int? skip = null,
      int? take = null)
   {
      if (predicate == null) throw new ArgumentNullException(nameof(predicate));

      IQueryable<T> query = _dbSet.Where(predicate);

      if (orderBy != null) query = orderBy(query);
      if (skip.HasValue) query = query.Skip(skip.Value);
      if (take.HasValue) query = query.Take(take.Value);

      return await query.ToListAsync();
   }

   public async Task<int> CountAsync(Expression<Func<T, bool>> predicate)
   {
      if (predicate == null) throw new ArgumentNullException(nameof(predicate));
    
      return await _dbSet.CountAsync(predicate);
   }
   
   public async Task<int> CountAsync()
   {
      return await _dbSet.CountAsync();
   }

   public async Task DeleteRange(IEnumerable<T> entities)
   {
      if (entities == null)
         throw new ArgumentNullException(nameof(entities));

      _dbSet.RemoveRange(entities);
      await _context.SaveChangesAsync();
   }

   public async Task ReloadAsync(T entity)
   {
      await _dbSet.Entry(entity).ReloadAsync();
   }
}