package sqlxopscope_test

import (
	"context"
	stderrs "errors"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sqlxopscope "github.com/infevocorp/goflexstore/sqlx/opscope"
)

func newDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE items (id INTEGER PRIMARY KEY AUTOINCREMENT, val TEXT)`)
	require.NoError(t, err)
	return db
}

func TestTx_NoTransaction_ReturnsDB(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ext := scope.Tx(context.Background())
	assert.NotNil(t, ext)
}

func TestBegin_End_Commit(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

	ctx, err := scope.Begin(ctx)
	require.NoError(t, err)

	_, err = scope.Tx(ctx).ExecContext(ctx, "INSERT INTO items (val) VALUES (?)", "hello")
	require.NoError(t, err)

	require.NoError(t, scope.End(ctx, nil))

	var count int
	require.NoError(t, db.Get(&count, "SELECT COUNT(*) FROM items"))
	assert.Equal(t, 1, count)
}

func TestBegin_End_Rollback(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

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
}

func TestBegin_Nested_IncrLevel(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

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
}

func TestEndWithRecover_PanicsAreRecovered(t *testing.T) {
	db := newDB(t)
	defer db.Close()

	scope := sqlxopscope.NewTransactionScope("test", db, nil)
	ctx := context.Background()

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
}
