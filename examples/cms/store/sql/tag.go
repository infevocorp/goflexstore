package sql

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/examples/cms/store/sql/dto"
	gormstore "github.com/jkaveri/goflexstore/gorm/store"
)

type TagStore struct {
	Store *gormstore.Store[*model.Tag, *dto.Tag, int64]
}
