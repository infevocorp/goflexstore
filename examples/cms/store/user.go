package store

import (
	"github.com/infevocorp/goflexstore/examples/cms/model"
	"github.com/infevocorp/goflexstore/store"
)

// UserStore is a store for users
type UserStore interface {
	store.Store[*model.User, int64]
}
