package store

import (
	"github.com/infevocorp/goflexstore/examples/cms/model"
	"github.com/infevocorp/goflexstore/store"
)

// ArticleStore is a store for articles
type ArticleStore interface {
	store.Store[*model.Article, int64]
}
