package model

import "time"

type Article struct {
	ID       int64
	Title    string
	Content  string
	AuthorID int64

	CreatedAt time.Time
	UpdatedAt time.Time

	Tags   []*Tag
	Author *User
}

func (a *Article) GetID() int64 {
	return a.ID
}
