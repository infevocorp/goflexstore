package dto

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u User) GetID() int64 {
	return u.ID
}
