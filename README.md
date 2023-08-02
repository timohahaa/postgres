### Postgres wrapper that I made, to reuse in later projects.
uses: 
- [github.com/jackc/pgx](https://github.com/jackc/pgx/)
- [github.com/Masterminds/squirrel](https://github.com/Masterminds/squirrel)

example:
```go
// create a postgres object
pg, err := postgres.New("postgres://myuser:mypassword@localhost:5432/db_name")

// you can also specify options:
pg, err := postgres.New("postgres://myuser:mypassword@localhost:5432/db_name", postgres.MaxConnPoolSize(5))

// here is a list of all options:
- ConnectionAttempts(int) - specify, how many times to try connect to postgres, if a connection fails
- ConnectionTimeout(time.Duration) - timeout between each connection attempt
- MaxConnPoolSize(int) - specify the size of the connection pool

// with NewOnce you can create a "singleton" postgres object - only one instance will be created, no matter how many times NewOnce was called
pgSingleton, err := postgres.NewOnce("postgres://myuser:mypassword@localhost")

// (options for NewOnce are also allowed)

// to query data:
// 1. build a query with squirrel sql-builder (optional)
sql, args, err := pg.Builder.
                    Select("*").
                    From("users").
                    ToSql()
// 2. ececute a query
rows, err := pg.ConnPool.Query(ctx, sql, args...)
// scan rows...
```
:)