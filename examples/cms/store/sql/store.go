package sql

import (
	"github.com/jkaveri/goflexstore/examples/cms/store"
	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
)

func NewStores(scope *gormopscope.TransactionScope) store.Stores {
	return store.Stores{
		Article: NewArticleStore(scope),
		User:    NewUserStore(scope),
	}
}
