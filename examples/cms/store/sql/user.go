package sql

import (
	"github.com/infevocorp/goflexstore/examples/cms/model"
	"github.com/infevocorp/goflexstore/examples/cms/store/sql/dto"
	gormopscope "github.com/infevocorp/goflexstore/gorm/opscope"
	gormstore "github.com/infevocorp/goflexstore/gorm/store"
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
