package model

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) GetID() int64 {
	return u.ID
}
