package store

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/store"
)

type TagStore interface {
	store.Store[*model.Tag, int64]
}
