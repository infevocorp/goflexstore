package sql

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/examples/cms/store/sql/dto"
	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	gormstore "github.com/jkaveri/goflexstore/gorm/store"
)

func NewArticleStore(scope *gormopscope.TransactionScope) *ArticleStore {
	return &ArticleStore{
		Store: gormstore.New[*model.Article, *dto.Article, int64](
			scope,
		),
	}
}

type ArticleStore struct {
	*gormstore.Store[*model.Article, *dto.Article, int64]
}
