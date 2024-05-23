package store

import (
	"github.com/infevocorp/goflexstore/examples/cms/model"
	"github.com/infevocorp/goflexstore/store"
)

type TagStore interface {
	store.Store[*model.Tag, int64]
}
