using System.Linq.Expressions;
using Webby.NotificationService.Data;

namespace Webby.NotificationService.Interfaces.Repositories;

public interface IRepository<TEntity> where TEntity : class
{
   Task<List<TEntity>> GetAll();
   Task<Guid> Add(TEntity entity);
   Task<TEntity?> FindById(Guid id);
   Task Update(TEntity item);
   Task<Guid> DeleteAsync(Guid id);
   Task<IEnumerable<TEntity>> GetByPredicate(Expression<Func<TEntity, bool>> predicate);
   Task<int> CountAsync();
   Task DeleteRange(IEnumerable<TEntity> entities);
   Task ReloadAsync(TEntity entity);

}