package sql

import (
	"gorm.io/gorm"

	"github.com/jkaveri/goflexstore/examples/cms/store/sql/dto"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		dto.Article{},
		dto.User{},
		dto.Tag{},
	)
}
