# goflexstore/sqlx [![Go Reference](https://pkg.go.dev/badge/github.com/infevocorp/goflexstore/sqlx.svg)](https://pkg.go.dev/github.com/infevocorp/goflexstore/sqlx)

`sqlx`-backed implementation of the [goflexstore](https://github.com/infevocorp/goflexstore) repository pattern for Go.

Supports **MySQL**, **PostgreSQL**, and **SQLite** via [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx).

## Installation

```bash
go get github.com/infevocorp/goflexstore@latest
go get github.com/infevocorp/goflexstore/sqlx@latest
```

## Packages

| Package | Import path | Purpose |
|---|---|---|
| `sqlxstore` | `goflexstore/sqlx/store` | Generic CRUD store |
| `sqlxopscope` | `goflexstore/sqlx/opscope` | Transaction management via context |
| `sqlxquery` | `goflexstore/sqlx/query` | SQL fragment builder from query params |
| `sqlxutils` | `goflexstore/sqlx/utils` | Reflection helpers (table name, field→column map) |

## Quick Start

### 1. Define your domain model and database row struct

```go
// Domain model — no db tags, no ORM coupling
type User struct {
    ID   int64
    Name string
    Age  int
}
func (u User) GetID() int64 { return u.ID }

// Database row struct — carries db: tags
type UserRow struct {
    ID   int64  `db:"id"`
    Name string `db:"name"`
    Age  int    `db:"age"`
}
func (u UserRow) GetID() int64 { return u.ID }
```

### 2. Create the store

```go
import (
    "github.com/jmoiron/sqlx"
    sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
    sqlxstore   "github.com/infevocorp/goflexstore/sqlx/store"
    sqlxquery   "github.com/infevocorp/goflexstore/sqlx/query"
)

db, _ := sqlx.Open("mysql", dsn)

opScope := sqlxopscope.NewTransactionScope("main", db, nil)

userStore := sqlxstore.New[User, UserRow, int64](
    opScope,
    sqlxstore.WithTable[User, UserRow, int64]("users"),
    sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectMySQL),
)
```

### 3. CRUD operations

```go
// Create
id, err := userStore.Create(ctx, User{Name: "Alice", Age: 30})

// Get
user, err := userStore.Get(ctx, query.Filter("ID", id))

// List with filters
users, err := userStore.List(ctx,
    query.Filter("Age", 18).WithOP(query.GTE),
    query.OrderBy("Name", false),
    query.Paginate(0, 20),
)

// Update all columns
err = userStore.Update(ctx, User{ID: id, Name: "Alice Updated", Age: 31})

// Partial update — only non-zero fields are written
err = userStore.PartialUpdate(ctx, User{ID: id, Name: "NewName"})

// Delete
err = userStore.Delete(ctx, query.Filter("ID", id))

// Upsert
id, err = userStore.Upsert(ctx, user, store.OnConflict{UpdateAll: true})
```

### 4. Transactions

```go
ctx, err = opScope.Begin(ctx)
defer opScope.EndWithRecover(ctx, &err) // commits on success, rolls back on error/panic

// ... multiple store operations sharing the same ctx
```

## Configuration Options

```go
sqlxstore.New[Entity, Row, ID](
    opScope,
    sqlxstore.WithTable[E, R, ID]("table_name"),           // override auto-derived table name
    sqlxstore.WithDialect[E, R, ID](sqlxquery.DialectPostgres),
    sqlxstore.WithPKColumn[E, R, ID]("uuid"),              // default: "id"
    sqlxstore.WithReturningID[E, R, ID](true),             // enable RETURNING for Postgres
    sqlxstore.WithBatchSize[E, R, ID](100),                // default: 50
    sqlxstore.WithConverter[E, R, ID](myConverter),        // custom Entity <-> Row converter
)
```

## Dialect Support

| Dialect constant | Database | Placeholder style |
|---|---|---|
| `sqlxquery.DialectMySQL` | MySQL (default) | `?` |
| `sqlxquery.DialectPostgres` | PostgreSQL | `$1`, `$2`, … |
| `sqlxquery.DialectSQLite` | SQLite | `?` |

For PostgreSQL, also set `WithReturningID(true)` to retrieve the inserted PK via `RETURNING`.

## Query Params

All store methods accept variadic `query.Param` values from `github.com/infevocorp/goflexstore/query`:

```go
query.Filter("FieldName", value)                    // WHERE field = ?
query.Filter("FieldName", value).WithOP(query.GT)   // WHERE field > ?
query.Filter("IDs", []int64{1, 2, 3})               // WHERE id IN (1,2,3)
query.OR(query.Filter("A", "x"), query.Filter("B", "y")) // WHERE (a = ? OR b = ?)
query.OrderBy("FieldName", true)                    // ORDER BY field DESC
query.Paginate(offset, limit)                       // LIMIT ? OFFSET ?
query.Select("Field1", "Field2")                    // SELECT field1, field2 (instead of *)
query.GroupBy("Field").WithHaving(query.Filter("Age", 18).WithOP(query.GT))
query.WithLock(query.LockTypeForUpdate)             // FOR UPDATE
query.WithHint("index_merge")                       // /*+ index_merge */
```

> `query.Preload` is **not** supported and returns `ErrPreloadNotSupported`.

## Table Name Inference

By default the table name is derived from the Row struct type:

- Strip trailing `Row` suffix → `UserRow` → `User`
- Convert to `snake_case` → `user`
- Pluralise → `users`

Override with `WithTable(...)` when the convention does not match.

## Transaction Scope

`TransactionScope` stores the active `*sqlx.Tx` in `context.Context` under a named key. Multiple scopes with different names can coexist (e.g., for multi-database setups). Nesting is supported — inner `Begin`/`End` calls increment/decrement a reference counter instead of opening a new transaction.

```go
// Read-only scope (READ COMMITTED, read-only isolation)
readScope := sqlxopscope.NewReadTransactionScope("read", db)

// Serializable write scope
writeScope := sqlxopscope.NewWriteTransactionScope("write", db)
```

## License

[MIT](../LICENSE.txt)
