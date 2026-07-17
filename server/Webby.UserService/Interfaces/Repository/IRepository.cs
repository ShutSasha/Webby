using System.Linq.Expressions;
using Webby.UserService.Data;

namespace Webby.UserService.Interfaces.Repository;

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
   Task<(List<TEntity> Items, int Total)> SearchAsync(
      string tableName,
      string columnName,
      string? searchText,
      int skip,
      int take,
      string? additionalWhere = null,
      object[]? parameters = null,
      string? orderBy = null,
      Func<AppDbContext, Expression<Func<TEntity, bool>>>? predicateFactory = null);
}