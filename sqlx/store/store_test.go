package sqlxstore_test

import (
	"context"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infevocorp/goflexstore/query"
	"github.com/infevocorp/goflexstore/store"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
	sqlxquery "github.com/infevocorp/goflexstore/sqlx/query"
	sqlxstore "github.com/infevocorp/goflexstore/sqlx/store"
)

// ---- fixtures ---------------------------------------------------------------

type UserRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (u UserRow) GetID() int64 { return u.ID }

type User struct {
	ID   int64
	Name string
	Age  int
}

func (u User) GetID() int64 { return u.ID }

// ---- helpers ----------------------------------------------------------------

func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE users (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT    NOT NULL,
		age  INTEGER NOT NULL DEFAULT 0
	)`)
	require.NoError(t, err)

	return db
}

func newTestStore(db *sqlx.DB) *sqlxstore.Store[User, UserRow, int64] {
	opScope := sqlxopscope.NewTransactionScope("test", db, nil)
	return sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
		sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectSQLite),
	)
}

// ---- Create -----------------------------------------------------------------

func TestCreate_AutoIncrement(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	id, err := s.Create(ctx, User{Name: "Alice", Age: 30})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestCreate_SequentialIDs(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	id1, err := s.Create(ctx, User{Name: "Alice", Age: 30})
	require.NoError(t, err)
	id2, err := s.Create(ctx, User{Name: "Bob", Age: 25})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id1)
	assert.Equal(t, int64(2), id2)
}

// ---- Get --------------------------------------------------------------------

func TestGet_Found(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, err := s.Create(ctx, User{Name: "Alice", Age: 30})
	require.NoError(t, err)

	user, err := s.Get(ctx, query.Filter("ID", int64(1)))
	require.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, 30, user.Age)
}

func TestGet_NotFound(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, err := s.Get(ctx, query.Filter("ID", int64(99)))
	assert.ErrorIs(t, err, store.ErrorNotFound)
}

func TestGet_Preload_NotSupported(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, err := s.Get(ctx, query.Preload("Something"))
	assert.ErrorIs(t, err, sqlxstore.ErrPreloadNotSupported)
}

// ---- List -------------------------------------------------------------------

func TestList_All(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})

	users, err := s.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestList_WithFilter(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})

	users, err := s.List(ctx, query.Filter("Age", 30))
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "Alice", users[0].Name)
}

func TestList_Paginate(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = s.Create(ctx, User{Name: "u", Age: i})
	}

	users, err := s.List(ctx, query.Paginate(2, 2))
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestList_OrderBy(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})
	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})

	users, err := s.List(ctx, query.OrderBy("Name", false))
	require.NoError(t, err)
	require.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
}

func TestList_InSlice(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})
	_, _ = s.Create(ctx, User{Name: "Carol", Age: 20})

	users, err := s.List(ctx, query.Filter("ID", []int64{1, 3}))
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

// ---- Count ------------------------------------------------------------------

func TestCount(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})

	n, err := s.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), n)
}

func TestCount_WithFilter(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})

	n, err := s.Count(ctx, query.Filter("Age", 30))
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)
}

// ---- Exists -----------------------------------------------------------------

func TestExists_True(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})

	ok, err := s.Exists(ctx, query.Filter("Name", "Alice"))
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestExists_False(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	ok, err := s.Exists(ctx, query.Filter("Name", "Ghost"))
	require.NoError(t, err)
	assert.False(t, ok)
}

// ---- Update -----------------------------------------------------------------

func TestUpdate_ByID(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	id, _ := s.Create(ctx, User{Name: "Alice", Age: 30})
	err := s.Update(ctx, User{ID: id, Name: "Alice Updated", Age: 31})
	require.NoError(t, err)

	user, _ := s.Get(ctx, query.Filter("ID", id))
	assert.Equal(t, "Alice Updated", user.Name)
	assert.Equal(t, 31, user.Age)
}

func TestUpdate_ByParam(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	err := s.Update(ctx, User{Name: "Renamed", Age: 30}, query.Filter("Name", "Alice"))
	require.NoError(t, err)

	users, _ := s.List(ctx)
	assert.Equal(t, "Renamed", users[0].Name)
}

func TestUpdate_NoIDNoParams_Error(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	err := s.Update(ctx, User{Name: "x"})
	assert.Error(t, err)
}

// ---- PartialUpdate ----------------------------------------------------------

func TestPartialUpdate_OnlyNonZero(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	id, _ := s.Create(ctx, User{Name: "Alice", Age: 30})
	// Age is zero → should not be updated
	err := s.PartialUpdate(ctx, User{ID: id, Name: "NewName"})
	require.NoError(t, err)

	user, _ := s.Get(ctx, query.Filter("ID", id))
	assert.Equal(t, "NewName", user.Name)
	assert.Equal(t, 30, user.Age) // unchanged
}

// ---- Delete -----------------------------------------------------------------

func TestDelete(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	_, _ = s.Create(ctx, User{Name: "Alice", Age: 30})
	_, _ = s.Create(ctx, User{Name: "Bob", Age: 25})

	err := s.Delete(ctx, query.Filter("Name", "Alice"))
	require.NoError(t, err)

	n, _ := s.Count(ctx)
	assert.Equal(t, int64(1), n)
}

// ---- CreateMany -------------------------------------------------------------

func TestCreateMany(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	users := []User{
		{Name: "A", Age: 1},
		{Name: "B", Age: 2},
		{Name: "C", Age: 3},
	}
	err := s.CreateMany(ctx, users)
	require.NoError(t, err)

	n, _ := s.Count(ctx)
	assert.Equal(t, int64(3), n)
}

// ---- Upsert -----------------------------------------------------------------

func TestUpsert_DoNothing_NoConflict(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	s := newTestStore(db)
	ctx := context.Background()

	id, err := s.Upsert(ctx, User{Name: "Alice", Age: 30}, store.OnConflict{DoNothing: true})
	require.NoError(t, err)
	assert.NotZero(t, id)

	n, _ := s.Count(ctx)
	assert.Equal(t, int64(1), n)
}

func TestUpsert_UpdateAll_Conflict(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	// Enable named PK for explicit upsert
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id   INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		age  INTEGER NOT NULL DEFAULT 0
	)`)
	require.NoError(t, err)

	opScope := sqlxopscope.NewTransactionScope("test", db, nil)
	s := sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("products"),
		sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectSQLite),
	)

	// Insert with explicit ID
	_, err = db.Exec("INSERT INTO products (id, name, age) VALUES (1, 'Alice', 30)")
	require.NoError(t, err)

	id, err := s.Upsert(ctx(t), User{ID: 1, Name: "Alice Updated", Age: 99}, store.OnConflict{UpdateAll: true})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)

	var row UserRow
	require.NoError(t, db.Get(&row, "SELECT id, name, age FROM products WHERE id = 1"))
	assert.Equal(t, "Alice Updated", row.Name)
	assert.Equal(t, 99, row.Age)
}

// ---- Transaction ------------------------------------------------------------

func TestTransaction_Commit(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	opScope := sqlxopscope.NewTransactionScope("tx", db, nil)
	s := sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
		sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectSQLite),
	)

	c := context.Background()
	c, err := opScope.Begin(c)
	require.NoError(t, err)

	_, err = s.Create(c, User{Name: "TxUser", Age: 10})
	require.NoError(t, err)

	require.NoError(t, opScope.End(c, nil))

	n, _ := s.Count(context.Background())
	assert.Equal(t, int64(1), n)
}

func TestTransaction_Rollback(t *testing.T) {
	db := newTestDB(t)
	defer db.Close()

	opScope := sqlxopscope.NewTransactionScope("tx", db, nil)
	s := sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
		sqlxstore.WithDialect[User, UserRow, int64](sqlxquery.DialectSQLite),
	)

	c := context.Background()
	c, err := opScope.Begin(c)
	require.NoError(t, err)

	_, err = s.Create(c, User{Name: "TxUser", Age: 10})
	require.NoError(t, err)

	_ = opScope.End(c, assert.AnError)

	n, _ := s.Count(context.Background())
	assert.Equal(t, int64(0), n)
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	return context.Background()
}
