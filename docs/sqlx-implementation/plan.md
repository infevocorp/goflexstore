# sqlx Implementation Plan

Implement `store.Store[T, ID]` backed by `jmoiron/sqlx`, mirroring the structure of the existing `gorm/` submodule.

---

## Repository layout

```
sqlx/                        ← new Go submodule (github.com/infevocorp/goflexstore/sqlx)
├── go.mod
├── store/
│   ├── store.go             ← Store[Entity, Row, ID] struct + New constructor
│   └── store_test.go
├── query/
│   ├── builder.go           ← QueryBuilder: translates query.Params → SQL fragments
│   ├── builder_test.go
│   └── placeholders.go      ← dialect-aware placeholder helper ($1/$2 vs ?)
├── opscope/
│   ├── txscope.go           ← TransactionScope (mirrors gorm/opscope)
│   └── txscope_test.go
└── utils/
    └── fieldtocolmap.go     ← reads `db:"column"` struct tags → map[fieldName]colName
```

---

## Core design decisions

### 1. Three generic type parameters — same as gormstore

```go
// Row is the sqlx-tagged struct (db tags), analogous to GORM's DTO.
// Entity is the clean domain model.
type Store[Entity store.Entity[ID], Row store.Entity[ID], ID comparable] struct {
    OpScope      *opscope.TransactionScope
    Converter    converter.Converter[Entity, Row, ID]
    QueryBuilder *query.Builder
    Table        string   // resolved once in New(), never at query time
    BatchSize    int
}
```

The existing `converter.NewReflect` from the root module works unchanged; no new converter is needed.

### 2. `New` constructor

```go
func New[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
    opScope *opscope.TransactionScope,
    options ...Option[Entity, Row, ID],
) *Store[Entity, Row, ID]
```

Default behaviour when no options are provided:
- `Converter` → `converter.NewReflect[Entity, Row, ID](nil)`
- `Table` → snake_case plural of the Row type name (same convention sqlx users expect), overridable via `WithTable("users")`
- `QueryBuilder` → `query.NewBuilder(query.WithFieldToColMap(utils.FieldToColMap(*new(Row))))`
- `BatchSize` → 50

### 3. Transaction scope

`opscope.TransactionScope` stores `*sqlx.Tx` in `context.Context` (keyed by scope name). API mirrors `gorm/opscope`:

```go
func (s *TransactionScope) Begin(ctx context.Context) (context.Context, error)
func (s *TransactionScope) End(ctx context.Context, err error) error
func (s *TransactionScope) EndWithRecover(ctx context.Context, errPtr *error)
func (s *TransactionScope) Tx(ctx context.Context) sqlx.ExtContext  // returns *sqlx.Tx or *sqlx.DB
```

`Tx` returns `sqlx.ExtContext` so callers work with either the transaction or the bare DB without branching.

### 4. Query builder

`query.Builder` walks `query.Params` and produces two outputs:

```go
type Result struct {
    Where   string   // "age = ? AND name != ?"  (placeholders already bound to dialect)
    Args    []any
    OrderBy string   // "name ASC, id DESC"
    Limit   int      // 0 = no limit
    Offset  int
    Cols    []string // for SELECT; nil = "*"
}

func (b *Builder) Build(params query.Params) Result
```

Supported param types and their SQL translation:

| `query.Param` type | SQL output |
|---|---|
| `FilterParam` (EQ/NEQ/GT/GTE/LT/LTE) | `col = ?`, `col != ?`, `col > ?`, etc. |
| `ORParam` | `(col1 = ? OR col2 = ?)` |
| `PaginateParam` | sets `Limit` / `Offset` |
| `OrderByParam` | appends to ORDER BY clause |
| `SelectParam` | sets `Cols` list |
| `GroupByParam` | appends GROUP BY + optional HAVING |
| `WithLockParam` (FOR UPDATE) | appends `FOR UPDATE` suffix |
| `WithHintParam` | prepended as comment hint `/*+ hint */` |
| `PreloadParam` | **not supported** — sqlx has no ORM-level preloading; return `ErrPreloadNotSupported` |

`placeholders.go` provides a `Rebind(dialect, sql)` wrapper around `sqlx.Rebind` so the builder can emit `?` placeholders internally and rebind to `$1`, `$2` for Postgres.

Field-name → column-name resolution uses `FieldToColMap` (reads `db:"colname"` tags), same algorithm as `gorm/utils/fieldtocolmap.go` but reading the `db` struct tag instead of the `gorm` tag.

### 5. SQL assembly per operation

Every method fetches the current `sqlx.ExtContext` from `OpScope.Tx(ctx)` then assembles and executes SQL directly — no ORM magic.

#### `Get`
```sql
SELECT {cols} FROM {table} WHERE {where} ORDER BY {orderby} LIMIT 1
```
Returns `store.ErrorNotFound` when `sqlx.ErrNoRows`.

#### `List`
```sql
SELECT {cols} FROM {table} [WHERE {where}] [ORDER BY {orderby}] [LIMIT n OFFSET m]
```

#### `Count`
```sql
SELECT COUNT(*) FROM {table} [WHERE {where}]
```

#### `Exists`
```sql
SELECT EXISTS(SELECT 1 FROM {table} [WHERE {where}])
```

#### `Create`
Reflects over the Row struct to build the column list and placeholder list, skipping the primary key if it is zero (auto-increment). Returns the new ID via `sql.Result.LastInsertId()`, or — for Postgres — uses a `RETURNING id` clause (opt-in via `WithReturningID(true)` option).

```sql
INSERT INTO {table} (col1, col2, ...) VALUES (?, ?, ...)
```

#### `CreateMany`
Batches rows in groups of `BatchSize`. Builds a multi-row VALUES clause per batch:

```sql
INSERT INTO {table} (col1, col2) VALUES (?, ?), (?, ?), ...
```

#### `Update`  (full — all columns including zero values)
```sql
UPDATE {table} SET col1 = ?, col2 = ?, ... WHERE {where}
```
If `params` is empty, falls back to `WHERE {pk_col} = ?` using the entity's ID. Returns an error if both are absent.

#### `PartialUpdate`  (non-zero fields only)
Uses reflection to collect fields with non-zero values, then builds a sparse SET clause:
```sql
UPDATE {table} SET col1 = ? [, col2 = ?] WHERE {where}
```

#### `Delete`
```sql
DELETE FROM {table} WHERE {where}
```

#### `Upsert`
Database-specific; three strategies controlled by `OnConflict`:

| `OnConflict` field | MySQL / SQLite | Postgres |
|---|---|---|
| `DoNothing: true` | `INSERT IGNORE` | `INSERT … ON CONFLICT DO NOTHING` |
| `UpdateAll: true` | `INSERT … ON DUPLICATE KEY UPDATE col=VALUES(col), …` | `INSERT … ON CONFLICT (cols) DO UPDATE SET col=EXCLUDED.col, …` |
| `Updates map[string]any` | SET literal values in the conflict clause | same |
| `UpdateColumns []string` | named columns only | named columns only |
| `OnConstraint string` | ignored (MySQL has no named constraints in upsert) | `ON CONFLICT ON CONSTRAINT name` |

The dialect is injected via `WithDialect(d Dialect)` option (enum: `DialectMySQL`, `DialectPostgres`, `DialectSQLite`).

---

## Package structure detail

### `sqlx/store/store.go`

```go
package sqlxstore

import (
    "context"
    "reflect"

    "github.com/jmoiron/sqlx"
    "github.com/infevocorp/goflexstore/converter"
    "github.com/infevocorp/goflexstore/query"
    "github.com/infevocorp/goflexstore/store"
    sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
    sqlxquery  "github.com/infevocorp/goflexstore/sqlx/query"
    sqlxutils  "github.com/infevocorp/goflexstore/sqlx/utils"
)

type Store[Entity store.Entity[ID], Row store.Entity[ID], ID comparable] struct {
    OpScope      *sqlxopscope.TransactionScope
    Converter    converter.Converter[Entity, Row, ID]
    QueryBuilder *sqlxquery.Builder
    Table        string
    BatchSize    int
}

// Implements store.Store[Entity, ID]
var _ store.Store[any, any] = (*Store[any, any, any])(nil) // compile-time check (adjust types)

func New[Entity store.Entity[ID], Row store.Entity[ID], ID comparable](
    opScope *sqlxopscope.TransactionScope,
    options ...Option[Entity, Row, ID],
) *Store[Entity, Row, ID] { … }
```

### `sqlx/query/builder.go`

```go
package sqlxquery

type Builder struct {
    FieldToColMap map[string]string
    Dialect       placeholders.Dialect
}

type Result struct {
    Where   string
    Args    []any
    OrderBy string
    Limit   int
    Offset  int
    Cols    []string
    Suffix  string // e.g. "FOR UPDATE"
}

func NewBuilder(opts ...Option) *Builder { … }
func (b *Builder) Build(params query.Params) Result { … }
func (b *Builder) getColName(name string) string { … }
```

### `sqlx/opscope/txscope.go`

```go
package sqlxopscope

type TransactionScope struct {
    Name      string
    DB        *sqlx.DB
    TxOptions *sql.TxOptions
}

func NewTransactionScope(name string, db *sqlx.DB, opts *sql.TxOptions) *TransactionScope
func (s *TransactionScope) Begin(ctx context.Context) (context.Context, error)
func (s *TransactionScope) End(ctx context.Context, err error) error
func (s *TransactionScope) EndWithRecover(ctx context.Context, errPtr *error)
func (s *TransactionScope) Tx(ctx context.Context) sqlx.ExtContext
```

### `sqlx/utils/fieldtocolmap.go`

```go
package sqlxutils

import "reflect"

// FieldToColMap reads `db:"colname"` struct tags.
// Returns map[StructFieldName]columnName.
func FieldToColMap(row any) map[string]string { … }
```

---

## `go.mod` for the new submodule

```
module github.com/infevocorp/goflexstore/sqlx

go 1.21

require (
    github.com/infevocorp/goflexstore v1.0.11
    github.com/jmoiron/sqlx v1.4.0
    github.com/pkg/errors v0.9.1
    github.com/stretchr/testify v1.8.4
)
```

Tests will add a SQLite driver (e.g. `github.com/glebarez/sqlite` / `github.com/glebarez/go-sqlite`) as a test-only dependency for in-process integration tests, same pattern used in `benchmark/`.

---

## Implementation order

1. **`sqlx/go.mod`** — module scaffold, no code yet.
2. **`sqlx/utils/fieldtocolmap.go`** + tests — reads `db` tags; trivial, unblocks everything else.
3. **`sqlx/query/placeholders.go`** — `Rebind` wrapper + `Dialect` type.
4. **`sqlx/query/builder.go`** + tests — translates every `query.Param` type; pure string manipulation, no DB needed.
5. **`sqlx/opscope/txscope.go`** + tests — transaction context management.
6. **`sqlx/store/store.go`** — implement all ten `store.Store` methods; integration-test each against in-memory SQLite.
7. Wire up a `compile_test.go` that asserts `*Store` satisfies `store.Store` at compile time.

---

## Key differences from `gorm/` submodule

| Concern | gorm/ | sqlx/ |
|---|---|---|
| SQL generation | GORM ORM | Hand-built strings in `query.Builder` |
| Column scan | GORM auto-scan | `sqlx.StructScan` / `sqlx.Select` using `db` tags |
| `Preload` | Supported (GORM joins) | Not supported — return `ErrPreloadNotSupported` |
| `WithHint` | `gorm.io/hints` clause | Prepended SQL comment `/*+ hint */` |
| `WithLock` | GORM `clause.Locking` | Appended `FOR UPDATE` string |
| Upsert dialect | GORM handles differences | Explicit per-dialect SQL in `Upsert` |
| Auto table name | GORM `TableName()` | `utils.TableName(row)` or `WithTable` option |
| Returning ID | GORM `Create` sets primary key | `LastInsertId()` or `RETURNING` clause (Postgres) |

---

## Scope boundaries (not in this implementation)

- No query caching or prepared statement pooling.
- No `Preload` / eager-loading — callers compose queries manually or use a service layer.
- No schema migration helpers — out of scope for a store layer.
- No soft-delete support — callers add a `deleted_at` filter themselves.
