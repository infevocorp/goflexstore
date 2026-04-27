package converter_test

import (
	"database/sql"
	"time"
)

type UserDTO struct {
	ID        int           `gorm:"column:id;primary_key"`
	Name      string        `gorm:"column:name"`
	Age       int           `gorm:"column:age"`
	IsAdmin   *sql.NullBool `gorm:"column:is_admin"`
	Disabled  sql.NullBool  `gorm:"column:disabled"`
	CreatedAt sql.NullTime  `gorm:"column:created_at"`

	Referer *UserDTO `gorm:"foreignKey:RefererID"`

	Friends []*UserDTO `gorm:"many2many:user_friends"`
}

func (d UserDTO) GetID() int {
	return d.ID
}

type User struct {
	ID        int
	Name      string
	Age       int
	Disabled  bool
	IsAdmin   bool
	CreatedAt time.Time

	Referer *User

	Friends []*User
}

func (e User) GetID() int {
	return e.ID
}
