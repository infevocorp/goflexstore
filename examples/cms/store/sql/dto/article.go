package dto

import (
	"database/sql"
)

type Article struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
	AuthorID int64  `gorm:"column:author_id"`

	CreatedAt sql.NullTime `gorm:"column:created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`

	Author *User `gorm:"foreignKey:AuthorID"`

	// Tags is the list of tags that this article has.
	//nolint:revive
	Tags []*Tag `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:ArticleID;references:ID;joinReferences:TagID"`
}

func (a Article) GetID() int64 {
	return a.ID
}
