package gormopscope_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormopscope "github.com/infevocorp/goflexstore/gorm/opscope"
)

func Test_NewWriteTransactionScope(t *testing.T) {
	// GIVEN
	var (
		name  = "test"
		db, _ = newTestDB(t)
	)

	// WHEN
	scope := gormopscope.NewWriteTransactionScope(name, db)

	// THEN
	require.NotNil(t, scope)
	assert.Equal(t, name, scope.Name)
	assert.NotNil(t, scope.RootTx)
	assert.Equal(t, &sql.TxOptions{Isolation: sql.LevelSerializable}, scope.TxOptions)
}

func Test_NewReadTransactionScope(t *testing.T) {
	// GIVEN
	var (
		name  = "test"
		db, _ = newTestDB(t)
	)

	// WHEN
	scope := gormopscope.NewReadTransactionScope(name, db)

	// THEN
	require.NotNil(t, scope)
	assert.Equal(t, name, scope.Name)
	assert.NotNil(t, scope.RootTx)
	assert.Equal(t, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}, scope.TxOptions)
}

func Test_NewTransactionScope(t *testing.T) {
	// GIVEN
	var (
		name      = "test"
		db, _     = newTestDB(t)
		txOptions = &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  true,
		}
	)

	// WHEN
	scope := gormopscope.NewTransactionScope(name, db, txOptions)

	// THEN
	require.NotNil(t, scope)
	assert.Equal(t, name, scope.Name)
	assert.NotNil(t, scope.RootTx)
	assert.Equal(t, txOptions, scope.TxOptions)
}

func Test_TransactionScope_Begin(t *testing.T) {
	t.Run("should-begin-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()

		// WHEN
		ctx2, err := scope.Begin(ctx)

		// THEN
		require.NoError(t, err)
		assert.NotEqual(t, ctx, ctx2)
	})

	t.Run("should-return-err", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		// WHEN
		ctx2, err := scope.Begin(ctx)

		// THEN
		assert.Error(t, err)
		assert.Equal(t, ctx, ctx2)
	})

	t.Run("double-begin-should-call-begin-once", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()

		// WHEN
		ctx2, err := scope.Begin(ctx)

		ctx3, err2 := scope.Begin(ctx2)

		// THEN
		assert.NoError(t, err)
		assert.NoError(t, err2)
		assert.NotEqual(t, ctx, ctx2)
		assert.Equal(t, ctx2, ctx3)
	})
}

func Test_TransactionScope_End(t *testing.T) {
	t.Run("should-do-nothing-if-begin-transaction-failed", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		ctx2, err := scope.Begin(ctx)

		// WHEN
		err = scope.End(ctx2, err)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should-do-nothing-if-not-in-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name  = "test"
			db, _ = newTestDB(t)
			scope = gormopscope.NewWriteTransactionScope(name, db)
			ctx   = context.Background()
		)

		// WHEN
		err := scope.End(ctx, nil)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should-do-nothing-if-having-upper-level-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		ctx3, err := scope.Begin(ctx2)
		require.NoError(t, err)

		// WHEN
		err = scope.End(ctx3, nil)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should-rollback-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		// WHEN
		err = scope.End(ctx2, assert.AnError)

		// THEN
		assert.Error(t, err)
	})

	t.Run("should-return-err-if-cannot-rollback-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback().WillReturnError(sql.ErrConnDone)

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		// WHEN
		err = scope.End(ctx2, assert.AnError)

		// THEN
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot rollback transaction")
	})

	t.Run("should-return-err-if-cannot-commit-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		// WHEN
		err = scope.End(ctx2, nil)

		// THEN
		assert.Error(t, err)
		assert.ErrorContains(t, err, "cannot commit transaction")
	})

	t.Run("should-commit-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		// WHEN
		err = scope.End(ctx2, nil)

		// THEN
		assert.NoError(t, err)
	})
}

func Test_TransactionScope_Tx(t *testing.T) {
	t.Run("should-return-tx-if-not-in-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name  = "test"
			db, _ = newTestDB(t)
			scope = gormopscope.NewWriteTransactionScope(name, db)
			ctx   = context.Background()
		)

		// WHEN
		tx := scope.Tx(ctx)

		// THEN
		assert.NotNil(t, tx)
	})

	t.Run("should-return-tx-if-in-transaction", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()

		ctx2, err := scope.Begin(ctx)
		require.NoError(t, err)

		// WHEN
		tx := scope.Tx(ctx2)

		// THEN
		assert.NotNil(t, tx)
		assert.NotEqual(t, db, tx)
	})
}

func Test_TransactionScope_EndWithRecover(t *testing.T) {
	t.Run("should-panic-if-err-pointer-is-nil", func(t *testing.T) {
		// GIVEN
		var (
			name  = "test"
			db, _ = newTestDB(t)
			scope = gormopscope.NewWriteTransactionScope(name, db)
			ctx   = context.Background()
		)

		// WHEN
		assert.Panics(t, func() {
			scope.EndWithRecover(ctx, nil)
		})
	})

	t.Run("should-recover-panic", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()

		var err error

		func() {
			ctx2, err2 := scope.Begin(ctx)
			require.NoError(t, err2)

			// WHEN
			defer scope.EndWithRecover(ctx2, &err)

			panic(assert.AnError)
		}()

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should-recover-panic-2", func(t *testing.T) {
		// GIVEN
		var (
			name        = "test"
			db, sqlMock = newTestDB(t)
			scope       = gormopscope.NewWriteTransactionScope(name, db)
			ctx         = context.Background()
		)

		sqlMock.ExpectBegin()
		sqlMock.ExpectRollback()

		var err error

		func() {
			ctx2, err2 := scope.Begin(ctx)
			require.NoError(t, err2)

			// WHEN
			defer scope.EndWithRecover(ctx2, &err)

			panic("test panic")
		}()

		assert.ErrorContains(t, err, "panic: test panic")
	})
}

func newTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.23"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})

	t.Cleanup(func() {
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	return gormDB, sqlMock
}
