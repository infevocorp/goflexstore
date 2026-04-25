package sqlxstore_test

import (
	"github.com/infevocorp/goflexstore/store"
	sqlxstore "github.com/infevocorp/goflexstore/sqlx/store"
)

// Compile-time assertion that *Store satisfies store.Store.
type compileCheckEntity struct{ id int64 }

func (e compileCheckEntity) GetID() int64 { return e.id }

var _ store.Store[compileCheckEntity, int64] = (*sqlxstore.Store[compileCheckEntity, compileCheckEntity, int64])(nil)
