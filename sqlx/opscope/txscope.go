// Package sqlxopscope manages sqlx transactions stored in context.Context.
// It mirrors the API of gorm/opscope but uses *sqlx.DB / *sqlx.Tx.
package sqlxopscope

import (
	"context"
	"database/sql"
	stderrs "errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var errBeginTx = stderrs.New("failed to begin transaction")

type (
	contextKey string

	scopeValue struct {
		tx    *sqlx.Tx
		level int16
	}
)

// TransactionScope wraps a *sqlx.DB and stores the active *sqlx.Tx in context.
type TransactionScope struct {
	Name      string
	DB        *sqlx.DB
	TxOptions *sql.TxOptions
}

// NewTransactionScope creates a TransactionScope with custom tx options.
func NewTransactionScope(name string, db *sqlx.DB, txOptions *sql.TxOptions) *TransactionScope {
	return &TransactionScope{
		Name:      name,
		DB:        db,
		TxOptions: txOptions,
	}
}

// NewWriteTransactionScope creates a TransactionScope with serializable isolation.
func NewWriteTransactionScope(name string, db *sqlx.DB) *TransactionScope {
	return NewTransactionScope(name, db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

// NewReadTransactionScope creates a read-only, read-committed TransactionScope.
func NewReadTransactionScope(name string, db *sqlx.DB) *TransactionScope {
	return NewTransactionScope(name, db, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
}

// Begin starts a new transaction and stores it in ctx, or increments the
// nesting level if a transaction is already active.
func (s *TransactionScope) Begin(ctx context.Context) (context.Context, error) {
	sv := s.getScopeValue(ctx)
	if sv != nil {
		sv.level++
		return ctx, nil
	}

	tx, err := s.DB.BeginTxx(ctx, s.TxOptions)
	if err != nil {
		return ctx, stderrs.Join(errBeginTx, err)
	}

	return s.setScopeValue(ctx, &scopeValue{tx: tx, level: 1}), nil
}

// End commits or rolls back the transaction. If the nesting level is greater
// than 1 it simply decrements the level and returns nil.
func (s *TransactionScope) End(ctx context.Context, err error) error {
	if stderrs.Is(err, errBeginTx) {
		return nil
	}

	sv := s.getScopeValue(ctx)
	if sv == nil {
		return nil
	}

	if sv.level > 1 {
		sv.level--
		return nil
	}

	if err != nil {
		if rb := sv.tx.Rollback(); rb != nil {
			return stderrs.Join(err, fmt.Errorf("rollback failed: %w", rb))
		}
		return err
	}

	if cm := sv.tx.Commit(); cm != nil {
		return fmt.Errorf("commit failed: %w", cm)
	}

	return nil
}

// EndWithRecover wraps End with panic recovery, writing the final error into *errPtr.
func (s *TransactionScope) EndWithRecover(ctx context.Context, errPtr *error) {
	if errPtr == nil {
		panic("errPtr cannot be nil")
	}

	err := *errPtr

	if r := recover(); r != nil {
		if ferr, ok := r.(error); ok {
			err = stderrs.Join(err, ferr)
		} else {
			err = stderrs.Join(err, fmt.Errorf("panic: %v", r))
		}
		*errPtr = err
	}

	if err2 := s.End(ctx, err); err2 != nil {
		*errPtr = stderrs.Join(err, err2)
	}
}

// Tx returns the active *sqlx.Tx from ctx, or the underlying *sqlx.DB when no
// transaction is in progress. Both satisfy sqlx.ExtContext.
func (s *TransactionScope) Tx(ctx context.Context) sqlx.ExtContext {
	if sv := s.getScopeValue(ctx); sv != nil {
		return sv.tx
	}
	return s.DB
}

func (s *TransactionScope) getScopeValue(ctx context.Context) *scopeValue {
	if v := ctx.Value(s.ctxKey()); v != nil {
		return v.(*scopeValue)
	}
	return nil
}

func (s *TransactionScope) setScopeValue(ctx context.Context, sv *scopeValue) context.Context {
	return context.WithValue(ctx, s.ctxKey(), sv)
}

func (s *TransactionScope) ctxKey() contextKey {
	return contextKey(s.Name)
}
