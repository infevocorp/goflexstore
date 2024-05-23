package sql

import (
	"github.com/infevocorp/goflexstore/examples/cms/model"
	"github.com/infevocorp/goflexstore/examples/cms/store/sql/dto"
	gormstore "github.com/infevocorp/goflexstore/gorm/store"
)

type TagStore struct {
	Store *gormstore.Store[*model.Tag, *dto.Tag, int64]
}
