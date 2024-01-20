package store

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/store"
)

// ArticleStore is a store for articles
type ArticleStore interface {
	store.Store[*model.Article, int64]
}
