package sql

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/examples/cms/store/sql/dto"
	gormstore "github.com/jkaveri/goflexstore/gorm/store"
)

type ArticleStore struct {
	Store *gormstore.Store[*model.Article, *dto.Article, int64]
}
