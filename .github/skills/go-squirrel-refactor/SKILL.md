---
name: go-squirrel-refactor
description: Use this skill when requested to migrate or refactor raw SQL queries to the 'squirrel' query builder in Go services using pgx.
---

## Execution Rules
1. **PostgreSQL Strict Requirement:** You MUST ALWAYS use `sq.StatementBuilder.PlaceholderFormat(sq.Dollar)` for any squirrel builder. Do not use default `?` placeholders.
2. **pgx Compatibility (NO RunWith):** NEVER use `.RunWith()`, `.ExecContext()`, `.QueryContext()`, or `.QueryRowContext()` from the squirrel package. The `pgx` library does not implement `squirrel.BaseRunner`.
3. **Execution Pattern (ToSql):** ALWAYS compile the query to a string and arguments array using `.ToSql()`. Then, execute the generated SQL natively using the provided `pgx` pool (`r.db`) or transaction (`tx`). 
   *Example:*
   `sql, args, err := query.ToSql()`
   `err = r.db.QueryRow(ctx, sql, args...).Scan(&id)`
4. **Imports:** Ensure `github.com/Masterminds/squirrel` (alias `sq`) is imported. Remove unused standard DB imports if necessary.
5. **Signatures:** NEVER change the input arguments or return types of the repository methods. Ensure existing `ctx context.Context` is passed as the first argument to `pgx` execution methods.
6. **Readability:** Format the squirrel builder chain cleanly (one method per line).
7. **EXISTS Queries & Select Restrictions:** - NEVER use `sq.Select(sq.Expr(...))`. The `.Select()` method only accepts strings, not `Sqlizer` interfaces.
   - To check for existence (EXISTS logic), build a query like `sq.Select("1").From(...).Where(...).Limit(1)`. Execute it using `QueryRow`, and handle `pgx.ErrNoRows` to return `false, nil`. If no error, return `true, nil`.

## Trigger:
Activate this skill when the user explicitly asks to "apply go-squirrel-refactor" or "migrate to squirrel".