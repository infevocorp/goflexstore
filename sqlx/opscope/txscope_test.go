package sqlxopscope_test

import (
	"context"
	stderrs "errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })
	return sqlx.NewDb(sqlDB, "mysql"), mock
}

func TestTx_NoTransaction_ReturnsDB(t *testing.T) {
	db, _ := newMockDB(t)
	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ext := scope.Tx(context.Background())
	assert.NotNil(t, ext)
}

func TestBegin_End_Commit(t *testing.T) {
	db, mock := newMockDB(t)
	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO items (val) VALUES (?)").
		WithArgs("hello").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT COUNT(*) FROM items").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ctx, err := scope.Begin(ctx)
	require.NoError(t, err)

	_, err = scope.Tx(ctx).ExecContext(ctx, "INSERT INTO items (val) VALUES (?)", "hello")
	require.NoError(t, err)

	require.NoError(t, scope.End(ctx, nil))

	var count int
	require.NoError(t, db.Get(&count, "SELECT COUNT(*) FROM items"))
	assert.Equal(t, 1, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBegin_End_Rollback(t *testing.T) {
	db, mock := newMockDB(t)
	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO items (val) VALUES (?)").
		WithArgs("hello").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	mock.ExpectQuery("SELECT COUNT(*) FROM items").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	ctx, err := scope.Begin(ctx)
	require.NoError(t, err)

	_, err = scope.Tx(ctx).ExecContext(ctx, "INSERT INTO items (val) VALUES (?)", "hello")
	require.NoError(t, err)

	testErr := stderrs.New("oops")
	endErr := scope.End(ctx, testErr)
	require.ErrorIs(t, endErr, testErr)

	var count int
	require.NoError(t, db.Get(&count, "SELECT COUNT(*) FROM items"))
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBegin_Nested_IncrLevel(t *testing.T) {
	db, mock := newMockDB(t)
	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO items (val) VALUES (?)").
		WithArgs("x").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT COUNT(*) FROM items").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ctx, err := scope.Begin(ctx)
	require.NoError(t, err)

	// Second begin should not start a new real transaction
	ctx, err = scope.Begin(ctx)
	require.NoError(t, err)

	_, err = scope.Tx(ctx).ExecContext(ctx, "INSERT INTO items (val) VALUES (?)", "x")
	require.NoError(t, err)

	// First End just decrements level
	require.NoError(t, scope.End(ctx, nil))

	// Second End commits
	require.NoError(t, scope.End(ctx, nil))

	var count int
	require.NoError(t, db.Get(&count, "SELECT COUNT(*) FROM items"))
	assert.Equal(t, 1, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEndWithRecover_PanicsAreRecovered(t *testing.T) {
	db, mock := newMockDB(t)
	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO items (val) VALUES (?)").
		WithArgs("y").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	mock.ExpectQuery("SELECT COUNT(*) FROM items").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	var err error
	func() {
		ctx, err = scope.Begin(ctx)
		require.NoError(t, err)
		defer scope.EndWithRecover(ctx, &err)

		_, _ = scope.Tx(ctx).ExecContext(ctx, "INSERT INTO items (val) VALUES (?)", "y")
		panic("something went wrong")
	}()

	assert.Error(t, err)

	var count int
	require.NoError(t, db.Get(&count, "SELECT COUNT(*) FROM items"))
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
