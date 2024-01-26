package gormopscope

import (
	"context"
	"database/sql"
	stderrs "errors"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var errBeginTx = errors.New("failed to begin transaction")

type (
	// contextKey is a string type used as a key in the context
	contextKey string

	// scopeValue contains the transaction and the transaction level
	// in the context
	scopeValue struct {
		tx    *gorm.DB
		level int16
	}
)

// NewWriteTransactionScope creates a new write transaction scope.
// This function initializes a TransactionScope with serializable isolation level, intended for write operations.
//
// Parameters:
//   - name: A string representing the name of the transaction scope, used as a context key.
//   - rootTx: The root *gorm.DB object to start a new session with specific configurations.
//
// Returns:
// A new TransactionScope object with write configuration.
//
// Example:
// Creating a write transaction scope:
//
//	writeScope := gormopscope.NewWriteTransactionScope("writeTx", rootTx)
//
// This example creates a new write transaction scope with serializable
// isolation level using the root transaction object 'rootTx'.
func NewWriteTransactionScope(name string, rootTx *gorm.DB) *TransactionScope {
	return NewTransactionScope(name, rootTx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

// NewReadTransactionScope creates a new read-only transaction scope.
// This function initializes a TransactionScope with read-committed isolation
// level and read-only mode, intended for read operations.
//
// Parameters:
//   - name: A string representing the name of the transaction scope, used as a context key.
//   - rootTx: The root *gorm.DB object to start a new session with specific configurations.
//
// Returns:
// A new TransactionScope object with read-only configuration.
//
// Example:
// Creating a read-only transaction scope:
//
//	readScope := gormopscope.NewReadTransactionScope("readTx", rootTx)
//
// This example creates a new read-only transaction scope with read-committed
// isolation level using the root transaction object 'rootTx'.
func NewReadTransactionScope(name string, rootTx *gorm.DB) *TransactionScope {
	return NewTransactionScope(name, rootTx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
}

// NewTransactionScope initializes a new transaction scope with specified settings.
//
// This function creates a TransactionScope, which serves as a wrapper for managing
// database transactions with specific options.
//
// Parameters:
//   - name: A string representing the name of the transaction scope, used as a key
//
// in the context.
//   - rootTx: The base *gorm.DB instance from which a new session will be started.
//     This session is configured to be a new DB connection, with default transactions
//     skipped and nested transactions disabled.
//   - txOptions: The transaction options specified as *sql.TxOptions. These options
//     define the isolation level and read-only status of the transaction.
//
// Returns:
// A pointer to the newly created TransactionScope instance.
//
// Example:
// Creating a new transaction scope for a write operation with serializable isolation:
//
//	rootDB := // obtain gorm.DB instance
//	txScope := gormopscope.NewTransactionScope(
//		"myWriteScope",
//		rootDB,
//		&sql.TxOptions{Isolation: sql.LevelSerializable},
//	)
//
// This example demonstrates how to create a new transaction scope named "myWriteScope"
// with serializable isolation level using a root gorm.DB instance.
func NewTransactionScope(name string, rootTx *gorm.DB, txOptions *sql.TxOptions) *TransactionScope {
	return &TransactionScope{
		Name: name,
		RootTx: rootTx.Session(&gorm.Session{
			NewDB:                    true,
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		}),
		TxOptions: txOptions,
	}
}

// TransactionScope represents a transaction context for database operations.
//
// The struct holds essential information for managing database transactions in a flexible and controlled manner.
//
// Fields:
//   - Name: A unique identifier for the transaction scope. This name is used as
//     a key in the context for managing nested transactions.
//   - RootTx: The root GORM database object (*gorm.DB) from which transaction
//     sessions are derived.
//   - TxOptions: Options for the transaction, including isolation level and
//     read-only status. It's a pointer to sql.TxOptions.
//
// Example:
// Creating a new TransactionScope for a read-write transaction:
//
//	rootTx := // initialize your *gorm.DB
//	txScope := gormopscope.NewTransactionScope("myTxScope", rootTx, &sql.TxOptions{
//		Isolation: sql.LevelSerializable,
//	})
//
// This example sets up a new transaction scope with serializable isolation level.
type TransactionScope struct {
	Name      string
	RootTx    *gorm.DB
	TxOptions *sql.TxOptions
}

// Begin starts a new transaction or increases the transaction level if already in a transaction.
// This method begins a new transaction scope using the RootTx and TxOptions.
// If the context already has an ongoing transaction, it increments the transaction
// level instead of starting a new one.
//
// Parameters:
//   - ctx: The current context.Context object.
//
// Returns:
//   - A new context.Context object containing the transaction scope.
//   - An error if beginning the transaction fails.
//
// Example:
// Starting a transaction scope:
//
//	ctx, err := txScope.Begin(context.Background())
//
// This example starts a new transaction scope or increments the transaction level if already in a transaction.
func (s *TransactionScope) Begin(ctx context.Context) (context.Context, error) {
	scopeVal := s.getScopeValue(ctx)

	if scopeVal != nil {
		scopeVal.level++
		return ctx, nil
	}

	tx := s.RootTx.WithContext(ctx).Begin(s.TxOptions)
	if tx.Error != nil {
		return ctx, stderrs.Join(errBeginTx, tx.Error)
	}

	scopeVal = &scopeValue{
		tx:    tx,
		level: 1,
	}

	return s.setScopeValue(ctx, scopeVal), nil
}

// End finalizes the transaction scope.
// This method ends the transaction scope by committing or rolling back the
// transaction. It decrements the transaction level if nested transactions exist.
// If an error is passed, it triggers a rollback.
//
// Parameters:
//   - ctx: The current context.Context object.
//   - err: An error encountered during the transaction, leading to a rollback.
//
// Returns:
//   - An error if committing or rolling back the transaction fails.
//
// Example:
// Ending a transaction scope:
//
//	err := txScope.End(ctx, someError)
//
// This example ends the transaction scope, committing if 'someError' is nil,
// or rolling back if 'someError' is non-nil.
func (s *TransactionScope) End(ctx context.Context, err error) error {
	if errors.Is(err, errBeginTx) {
		return nil
	}

	scopeVal := s.getScopeValue(ctx)
	if scopeVal == nil {
		return nil
	}

	if scopeVal.level > 1 {
		scopeVal.level--
		return nil
	}

	if err != nil {
		if err2 := scopeVal.tx.Rollback().Error; err2 != nil {
			return stderrs.Join(err, errors.Wrap(err2, "cannot rollback transaction"))
		}

		return err
	}

	if err := scopeVal.tx.Commit().Error; err != nil {
		return errors.Wrap(err, "cannot commit transaction")
	}

	return nil
}

// Tx retrieves the current transaction from the context, if available, or otherwise returns the root transaction.
//
// This function checks for an active transaction associated with the current context. If such a transaction exists,
// it returns this transaction. Otherwise, it falls back to the root transaction initially set in the transaction scope.
//
// Parameters:
//   - ctx: A context.Context instance which may contain an ongoing transaction.
//
// Returns:
//   - *gorm.DB: The current transaction if present in the context; otherwise, the root transaction.
//
// Example:
// Working with transactions in a context:
//
//	func someDatabaseOperation(ctx context.Context, scope *TransactionScope) error {
//		tx := scope.Tx(ctx)
//		// Perform operations using tx...
//	}
//
// This example demonstrates retrieving the current transaction from the context using the Tx method. If no
// transaction is found in the context, the root transaction of the scope is used for database operations.
func (s *TransactionScope) Tx(ctx context.Context) *gorm.DB {
	sv := s.getScopeValue(ctx)
	if sv != nil {
		return sv.tx
	}

	return s.RootTx
}

// EndWithRecover implements the OperationScope interface by ending the transaction scope
// with a recovered error. It ensures that the transaction is correctly closed in the event of a panic.
//
// This method is crucial for handling unexpected errors and panics during a transaction. It
// attempts to recover from a panic, join any existing errors with the recovered error, and then
// properly end the transaction scope by either committing or rolling back, depending on the error state.
//
// Parameters:
//   - ctx: The context in which the transaction is operating. It is used for passing the transaction scope.
//   - errPtr: A pointer to an error variable that will be updated with the final error
//     state after recovery and transaction closure.
//
// If a panic occurs, the panic value is converted to an error (if it's not already an error) and combined
// with the existing error pointed to by errPtr. The transaction is then ended, and any error from ending
// the transaction is also combined with the existing error.
//
// It is important to pass a non-nil errPtr, as a nil pointer will result in a panic.
//
// Example:
//
//	var err error
//	ctx, err := transactionScope.Begin(context.Background())
//	if err != nil {
//	    // handle error
//	}
//	defer transactionScope.EndWithRecover(ctx, &err)
//	// perform operations within the transaction
//
// This example demonstrates the use of EndWithRecover in a typical transaction workflow.
// The transaction is begun, and then operations are performed within the transaction scope. If a panic
// occurs, EndWithRecover ensures the transaction is closed properly, and any errors are captured.
func (s *TransactionScope) EndWithRecover(ctx context.Context, errPtr *error) {
	if errPtr == nil {
		panic("err pointer cannot be nil")
	}

	err := *errPtr

	if r := recover(); r != nil {
		if ferr, ok := r.(error); ok {
			err = stderrs.Join(err, ferr)
		} else {
			err = stderrs.Join(err, errors.Errorf("panic: %v", r))
		}

		*errPtr = err
	}

	if err2 := s.End(ctx, err); err2 != nil {
		*errPtr = stderrs.Join(err, err2)
	}
}

func (s *TransactionScope) getScopeValue(ctx context.Context) *scopeValue {
	if val := ctx.Value(s.getCtxKey()); val != nil {
		return val.(*scopeValue)
	}

	return nil
}

func (s *TransactionScope) setScopeValue(ctx context.Context, scopeVal *scopeValue) context.Context {
	return context.WithValue(ctx, s.getCtxKey(), scopeVal)
}

func (s *TransactionScope) getCtxKey() contextKey {
	return contextKey(s.Name)
}
