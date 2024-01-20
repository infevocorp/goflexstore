package model

type Tag struct {
	ID       int64
	Slug     string
	Articles []*Article
}

func (t *Tag) GetID() int64 {
	return t.ID
}
