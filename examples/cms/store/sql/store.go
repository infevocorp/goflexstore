package sql

import (
	"github.com/infevocorp/goflexstore/examples/cms/store"
	gormopscope "github.com/infevocorp/goflexstore/gorm/opscope"
)

func NewStores(scope *gormopscope.TransactionScope) store.Stores {
	return store.Stores{
		Article: NewArticleStore(scope),
		User:    NewUserStore(scope),
	}
}
