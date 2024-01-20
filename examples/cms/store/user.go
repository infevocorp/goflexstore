package store

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/store"
)

// UserStore is a store for users
type UserStore interface {
	store.Store[*model.User, int64]
}
