package dto

import (
	"database/sql"
)

type Article struct {
	ID        int64        `gorm:"column:id;primaryKey;uuid"`
	Title     string       `gorm:"column:title"`
	Content   string       `gorm:"column:content"`
	Tags      []*Tag       `gorm:"column:tags"`
	AuthorID  int64        `gorm:"column:author_id"`
	Author    *User        `gorm:"foreignkey:AuthorID"`
	CreatedAt sql.NullTime `gorm:"column:created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
}

func (a Article) GetID() int64 {
	return a.ID
}
