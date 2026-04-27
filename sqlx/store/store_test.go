package sqlxstore_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/infevocorp/goflexstore/query"
	"github.com/infevocorp/goflexstore/store"
	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
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

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })
	return sqlx.NewDb(sqlDB, "mysql"), mock
}

func newTestStore(db *sqlx.DB) *sqlxstore.Store[User, UserRow, int64] {
	opScope := sqlxopscope.NewTransactionScope("test", db, nil)
	return sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
	)
}

// ---- Create -----------------------------------------------------------------

func TestCreate_AutoIncrement(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?)").
		WithArgs("Alice", 30).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := s.Create(context.Background(), User{Name: "Alice", Age: 30})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_SequentialIDs(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?)").
		WithArgs("Alice", 30).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?)").
		WithArgs("Bob", 25).
		WillReturnResult(sqlmock.NewResult(2, 1))

	id1, err := s.Create(context.Background(), User{Name: "Alice", Age: 30})
	require.NoError(t, err)
	id2, err := s.Create(context.Background(), User{Name: "Bob", Age: 25})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id1)
	assert.Equal(t, int64(2), id2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Get --------------------------------------------------------------------

func TestGet_Found(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(1), "Alice", 30)
	mock.ExpectQuery("SELECT * FROM users WHERE id = ? LIMIT 1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	user, err := s.Get(context.Background(), query.Filter("ID", int64(1)))
	require.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, 30, user.Age)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectQuery("SELECT * FROM users WHERE id = ? LIMIT 1").
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}))

	_, err := s.Get(context.Background(), query.Filter("ID", int64(99)))
	assert.ErrorIs(t, err, store.ErrorNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_Preload_NotSupported(t *testing.T) {
	db, _ := newMockDB(t)
	s := newTestStore(db)

	_, err := s.Get(context.Background(), query.Preload("Something"))
	assert.ErrorIs(t, err, sqlxstore.ErrPreloadNotSupported)
}

// ---- List -------------------------------------------------------------------

func TestList_All(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(1), "Alice", 30).
		AddRow(int64(2), "Bob", 25)
	mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)

	users, err := s.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilter(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(1), "Alice", 30)
	mock.ExpectQuery("SELECT * FROM users WHERE age = ?").
		WithArgs(30).
		WillReturnRows(rows)

	users, err := s.List(context.Background(), query.Filter("Age", 30))
	require.NoError(t, err)
	require.Len(t, users, 1)
	assert.Equal(t, "Alice", users[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_Paginate(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(3), "u", 2).
		AddRow(int64(4), "u", 3)
	mock.ExpectQuery("SELECT * FROM users LIMIT 2 OFFSET 2").WillReturnRows(rows)

	users, err := s.List(context.Background(), query.Paginate(2, 2))
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_OrderBy(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(2), "Alice", 30).
		AddRow(int64(1), "Bob", 25)
	mock.ExpectQuery("SELECT * FROM users ORDER BY name ASC").WillReturnRows(rows)

	users, err := s.List(context.Background(), query.OrderBy("Name", false))
	require.NoError(t, err)
	require.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_InSlice(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(int64(1), "Alice", 30).
		AddRow(int64(3), "Carol", 20)
	mock.ExpectQuery("SELECT * FROM users WHERE id IN (?, ?)").
		WithArgs(int64(1), int64(3)).
		WillReturnRows(rows)

	users, err := s.List(context.Background(), query.Filter("ID", []int64{1, 3}))
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Count ------------------------------------------------------------------

func TestCount(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectQuery("SELECT COUNT(*) FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(2)))

	n, err := s.Count(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCount_WithFilter(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectQuery("SELECT COUNT(*) FROM users WHERE age = ?").
		WithArgs(30).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

	n, err := s.Count(context.Background(), query.Filter("Age", 30))
	require.NoError(t, err)
	assert.Equal(t, int64(1), n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Exists -----------------------------------------------------------------

func TestExists_True(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)").
		WithArgs("Alice").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ok, err := s.Exists(context.Background(), query.Filter("Name", "Alice"))
	require.NoError(t, err)
	assert.True(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExists_False(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)").
		WithArgs("Ghost").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	ok, err := s.Exists(context.Background(), query.Filter("Name", "Ghost"))
	require.NoError(t, err)
	assert.False(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Update -----------------------------------------------------------------

func TestUpdate_ByID(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("UPDATE users SET name = ?, age = ? WHERE id = ?").
		WithArgs("Alice Updated", 31, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.Update(context.Background(), User{ID: 1, Name: "Alice Updated", Age: 31})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_ByParam(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("UPDATE users SET name = ?, age = ? WHERE name = ?").
		WithArgs("Renamed", 30, "Alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.Update(context.Background(), User{Name: "Renamed", Age: 30}, query.Filter("Name", "Alice"))
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_NoIDNoParams_Error(t *testing.T) {
	db, _ := newMockDB(t)
	s := newTestStore(db)

	err := s.Update(context.Background(), User{Name: "x"})
	assert.Error(t, err)
}

// ---- PartialUpdate ----------------------------------------------------------

func TestPartialUpdate_OnlyNonZero(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("UPDATE users SET name = ? WHERE id = ?").
		WithArgs("NewName", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.PartialUpdate(context.Background(), User{ID: 1, Name: "NewName"})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Delete -----------------------------------------------------------------

func TestDelete(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("DELETE FROM users WHERE name = ?").
		WithArgs("Alice").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.Delete(context.Background(), query.Filter("Name", "Alice"))
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- CreateMany -------------------------------------------------------------

func TestCreateMany(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?), (?, ?), (?, ?)").
		WithArgs("A", 1, "B", 2, "C", 3).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := s.CreateMany(context.Background(), []User{
		{Name: "A", Age: 1},
		{Name: "B", Age: 2},
		{Name: "C", Age: 3},
	})
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Upsert -----------------------------------------------------------------

func TestUpsert_DoNothing_NoConflict(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("INSERT IGNORE INTO users (name, age) VALUES (?, ?)").
		WithArgs("Alice", 30).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := s.Upsert(context.Background(), User{Name: "Alice", Age: 30}, store.OnConflict{DoNothing: true})
	require.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsert_UpdateAll_Conflict(t *testing.T) {
	db, mock := newMockDB(t)
	s := newTestStore(db)

	mock.ExpectExec("INSERT INTO users (id, name, age) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name), age = VALUES(age)").
		WithArgs(int64(1), "Alice Updated", 99).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := s.Upsert(context.Background(), User{ID: 1, Name: "Alice Updated", Age: 99}, store.OnConflict{UpdateAll: true})
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ---- Transaction ------------------------------------------------------------

func TestTransaction_Commit(t *testing.T) {
	db, mock := newMockDB(t)
	opScope := sqlxopscope.NewTransactionScope("tx", db, nil)
	s := sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
	)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?)").
		WithArgs("TxUser", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT COUNT(*) FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(1)))

	c := context.Background()
	c, err := opScope.Begin(c)
	require.NoError(t, err)

	_, err = s.Create(c, User{Name: "TxUser", Age: 10})
	require.NoError(t, err)

	require.NoError(t, opScope.End(c, nil))

	n, _ := s.Count(context.Background())
	assert.Equal(t, int64(1), n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransaction_Rollback(t *testing.T) {
	db, mock := newMockDB(t)
	opScope := sqlxopscope.NewTransactionScope("tx", db, nil)
	s := sqlxstore.New[User, UserRow, int64](
		opScope,
		sqlxstore.WithTable[User, UserRow, int64]("users"),
	)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users (name, age) VALUES (?, ?)").
		WithArgs("TxUser", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	mock.ExpectQuery("SELECT COUNT(*) FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(0)))

	c := context.Background()
	c, err := opScope.Begin(c)
	require.NoError(t, err)

	_, err = s.Create(c, User{Name: "TxUser", Age: 10})
	require.NoError(t, err)

	_ = opScope.End(c, assert.AnError)

	n, _ := s.Count(context.Background())
	assert.Equal(t, int64(0), n)
	assert.NoError(t, mock.ExpectationsWereMet())
}
