package sql

import (
	"github.com/jkaveri/goflexstore/examples/cms/model"
	"github.com/jkaveri/goflexstore/examples/cms/store/sql/dto"
	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	gormstore "github.com/jkaveri/goflexstore/gorm/store"
)

type UserStore struct {
	*gormstore.Store[*model.User, *dto.User, int64]
}

func NewUserStore(scope *gormopscope.TransactionScope) *UserStore {
	return &UserStore{
		Store: gormstore.New[*model.User, *dto.User, int64](
			scope,
		),
	}
}
