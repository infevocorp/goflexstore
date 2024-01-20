package gormstore_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

type UserDTO struct {
	ID       int           `gorm:"column:id;primary_key"`
	Name     string        `gorm:"column:name"`
	Age      int           `gorm:"column:age"`
	IsAdmin  *sql.NullBool `gorm:"column:is_admin"`
	Disabled sql.NullBool  `gorm:"column:disabled"`
}

func (d UserDTO) GetID() int {
	return d.ID
}

type User struct {
	ID       int
	Name     string
	Age      int
	Disabled bool
	IsAdmin  bool
}

func (e User) GetID() int {
	return e.ID
}

func newTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.23"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})

	t.Cleanup(func() {
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	return gormDB, sqlMock
}
