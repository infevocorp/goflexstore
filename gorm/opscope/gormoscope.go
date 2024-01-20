package gormopscope

import (
	"context"
	"database/sql"
	stderrs "errors"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var errBeginTx = errors.New("failed to begin transaction")

type contextKey string

type scopeValue struct {
	tx    *gorm.DB
	level int16
}

// NewWriteTransactionScope creates new write transaction scope
func NewWriteTransactionScope(name string, rootTx *gorm.DB) *TransactionScope {
	return NewTransactionScope(name, rootTx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

// NewReadTransactionScope creates new read only transaction scope
func NewReadTransactionScope(name string, rootTx *gorm.DB) *TransactionScope {
	return NewTransactionScope(name, rootTx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
}

// NewTransactionScope creates new transaction scope
//
// `name` is the name of the transaction scope, it will be used as context key
//
// `rootTx` root *gorm.DB to start new session with configuration: NewDB, SkipDefaultTransaction, DisableNestedTransaction
//
// `txOptions` is the transaction options
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

// TransactionScope is a struct that holds the root transaction and the transaction options
type TransactionScope struct {
	Name      string
	RootTx    *gorm.DB
	TxOptions *sql.TxOptions
}

// Begin begins the transaction scope
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

// End ends the transaction scope
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

// Tx returns current transaction in context if exists, otherwise returns root transaction
func (s *TransactionScope) Tx(ctx context.Context) *gorm.DB {
	sv := s.getScopeValue(ctx)
	if sv != nil {
		return sv.tx
	}

	return s.RootTx
}

// EndWithRecover ends the transaction scope with recovered error
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
