package bench_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	gormopscope "github.com/infevocorp/goflexstore/gorm/opscope"
	gormstore "github.com/infevocorp/goflexstore/gorm/store"
	"github.com/infevocorp/goflexstore/query"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	sqlxstore "github.com/infevocorp/goflexstore/sqlx/store"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User is the clean domain entity with no ORM tags.
type User struct {
	ID    int64
	Name  string
	Email string
	Age   int
}

func (u User) GetID() int64 { return u.ID }

// UserDTO is the GORM data-transfer object mapped to the "users" table.
type UserDTO struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `gorm:"column:name"`
	Email string `gorm:"column:email"`
	Age   int    `gorm:"column:age"`
}

func (d UserDTO) GetID() int64        { return d.ID }
func (UserDTO) TableName() string     { return "users" }

// UserRow is the sqlx scan target and implements store.Entity[int64].
type UserRow struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

func (u UserRow) GetID() int64 { return u.ID }

// ---- DB setup helpers ----

func newGormDB(b testing.TB) *gorm.DB {
	b.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		b.Fatalf("gorm.Open: %v", err)
	}
	if err := db.AutoMigrate(&UserDTO{}); err != nil {
		b.Fatalf("AutoMigrate: %v", err)
	}
	return db
}

func newSqlxDB(b testing.TB) *sqlx.DB {
	b.Helper()
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		b.Fatalf("sqlx.Open: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		name  TEXT    NOT NULL,
		email TEXT    NOT NULL,
		age   INTEGER NOT NULL
	)`)
	if err != nil {
		b.Fatalf("create table: %v", err)
	}
	return db
}

func newGormStore(db *gorm.DB) *gormstore.Store[User, UserDTO, int64] {
	opScope := gormopscope.NewTransactionScope("bench", db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	return gormstore.New[User, UserDTO, int64](opScope)
}

// ---- Seed helpers ----

func seedGorm(b testing.TB, st *gormstore.Store[User, UserDTO, int64], n int) []int64 {
	b.Helper()
	ctx := context.Background()
	ids := make([]int64, n)
	for i := 0; i < n; i++ {
		id, err := st.Create(ctx, User{
			Name:  fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + (i % 50),
		})
		if err != nil {
			b.Fatalf("seedGorm: %v", err)
		}
		ids[i] = id
	}
	return ids
}

func seedSqlx(b testing.TB, db *sqlx.DB, n int) []int64 {
	b.Helper()
	ids := make([]int64, n)
	for i := 0; i < n; i++ {
		res, err := db.Exec(
			`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("user%d@example.com", i),
			20+(i%50),
		)
		if err != nil {
			b.Fatalf("seedSqlx: %v", err)
		}
		id, _ := res.LastInsertId()
		ids[i] = id
	}
	return ids
}

// ---- Get by ID ----

func BenchmarkGormStore_Get(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	ids := seedGorm(b, st, 1)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.Get(ctx, query.Filter("ID", ids[0])); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_Get(b *testing.B) {
	db := newSqlxDB(b)
	ids := seedSqlx(b, db, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var u UserRow
		if err := db.Get(&u, `SELECT id, name, email, age FROM users WHERE id = ?`, ids[0]); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- List with pagination ----

func BenchmarkGormStore_List(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	seedGorm(b, st, 100)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.List(ctx, query.Paginate(0, 20)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_List(b *testing.B) {
	db := newSqlxDB(b)
	seedSqlx(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []UserRow
		if err := db.Select(&users, `SELECT id, name, email, age FROM users LIMIT 20 OFFSET 0`); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- List with filter + order + pagination ----

func BenchmarkGormStore_ListWithFilter(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	seedGorm(b, st, 100)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.List(ctx,
			query.Filter("Age", 30),
			query.OrderBy("Name", false),
			query.Paginate(0, 20),
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_ListWithFilter(b *testing.B) {
	db := newSqlxDB(b)
	seedSqlx(b, db, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []UserRow
		if err := db.Select(&users,
			`SELECT id, name, email, age FROM users WHERE age = ? ORDER BY name LIMIT 20 OFFSET 0`,
			30,
		); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- Create ----

func BenchmarkGormStore_Create(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.Create(ctx, User{
			Name:  fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   25,
		}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_Create(b *testing.B) {
	db := newSqlxDB(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := db.Exec(
			`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("user%d@example.com", i),
			25,
		); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- Update ----

func BenchmarkGormStore_Update(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	ids := seedGorm(b, st, 1)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := st.Update(ctx,
			User{ID: ids[0], Name: fmt.Sprintf("updated%d", i), Email: "updated@example.com", Age: 30},
			query.Filter("ID", ids[0]),
		); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_Update(b *testing.B) {
	db := newSqlxDB(b)
	ids := seedSqlx(b, db, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := db.Exec(
			`UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?`,
			fmt.Sprintf("updated%d", i), "updated@example.com", 30, ids[0],
		); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- Delete ----
// Pre-seeds b.N records so each iteration deletes a unique row,
// keeping the timer clean (seeding happens before b.ResetTimer).

func BenchmarkGormStore_Delete(b *testing.B) {
	db := newGormDB(b)
	st := newGormStore(db)
	ids := seedGorm(b, st, b.N)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := st.Delete(ctx, query.Filter("ID", ids[i])); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSqlx_Delete(b *testing.B) {
	db := newSqlxDB(b)
	ids := seedSqlx(b, db, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := db.Exec(`DELETE FROM users WHERE id = ?`, ids[i]); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore helpers ----

func newSqlxStore(db *sqlx.DB) *sqlxstore.Store[User, UserRow, int64] {
	opScope := sqlxopscope.NewTransactionScope("bench", db, nil)
	return sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
		sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectSQLite),
	)
}

func seedSqlxStore(b testing.TB, st *sqlxstore.Store[User, UserRow, int64], n int) []int64 {
	b.Helper()
	ctx := context.Background()
	ids := make([]int64, n)
	for i := 0; i < n; i++ {
		id, err := st.Create(ctx, User{
			Name:  fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + (i % 50),
		})
		if err != nil {
			b.Fatalf("seedSqlxStore: %v", err)
		}
		ids[i] = id
	}
	return ids
}

// ---- SqlxStore: Get by ID ----

func BenchmarkSqlxStore_Get(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	ids := seedSqlxStore(b, st, 1)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.Get(ctx, query.Filter("ID", ids[0])); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore: List with pagination ----

func BenchmarkSqlxStore_List(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	seedSqlxStore(b, st, 100)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.List(ctx, query.Paginate(0, 20)); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore: List with filter + order + pagination ----

func BenchmarkSqlxStore_ListWithFilter(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	seedSqlxStore(b, st, 100)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.List(ctx,
			query.Filter("Age", 30),
			query.OrderBy("Name", false),
			query.Paginate(0, 20),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore: Create ----

func BenchmarkSqlxStore_Create(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := st.Create(ctx, User{
			Name:  fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   25,
		}); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore: Update ----

func BenchmarkSqlxStore_Update(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	ids := seedSqlxStore(b, st, 1)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := st.Update(ctx,
			User{ID: ids[0], Name: fmt.Sprintf("updated%d", i), Email: "updated@example.com", Age: 30},
			query.Filter("ID", ids[0]),
		); err != nil {
			b.Fatal(err)
		}
	}
}

// ---- SqlxStore: Delete ----
// Pre-seeds b.N records so each iteration deletes a unique row.

func BenchmarkSqlxStore_Delete(b *testing.B) {
	db := newSqlxDB(b)
	st := newSqlxStore(db)
	ids := seedSqlxStore(b, st, b.N)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := st.Delete(ctx, query.Filter("ID", ids[i])); err != nil {
			b.Fatal(err)
		}
	}
}
