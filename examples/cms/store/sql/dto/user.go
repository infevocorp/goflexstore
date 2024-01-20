package dto

import "database/sql"

type User struct {
	ID        int64        `gorm:"column:id;primary_key"`
	Name      string       `gorm:"column:name"`
	Email     string       `gorm:"column:email"`
	CreatedAt sql.NullTime `gorm:"column:created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
}

func (u User) GetID() int64 {
	return u.ID
}
