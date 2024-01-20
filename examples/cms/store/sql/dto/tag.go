package dto

type Tag struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement"`

	Slug string `gorm:"column:slug"`

	// Articles is the list of articles that have this tag.
	//nolint:revive
	Articles []*Article `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:TagID;references:ID;joinReferences:ArticleID"`
}

func (t *Tag) GetID() int64 {
	return t.ID
}
