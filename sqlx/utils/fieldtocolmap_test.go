package sqlxutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sqlxutils "github.com/infevocorp/goflexstore/sqlx/utils"
)

func TestFieldToColMap(t *testing.T) {
	type UserRow struct {
		ID        int64  `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		ignored   string //nolint:unused
	}

	m := sqlxutils.FieldToColMap(UserRow{})
	assert.Equal(t, map[string]string{
		"ID":        "id",
		"FirstName": "first_name",
		"LastName":  "last_name",
	}, m)
}

func TestFieldToColMap_NoTag(t *testing.T) {
	type Simple struct {
		ID   int64
		Name string
	}

	m := sqlxutils.FieldToColMap(Simple{})
	assert.Equal(t, map[string]string{
		"ID":   "ID",
		"Name": "Name",
	}, m)
}

func TestFieldToColMap_OmitEmpty(t *testing.T) {
	type Row struct {
		ID int64 `db:"id,omitempty"`
	}

	m := sqlxutils.FieldToColMap(Row{})
	assert.Equal(t, "id", m["ID"])
}

func TestFieldToColMap_DashSkipped(t *testing.T) {
	type Row struct {
		ID     int64  `db:"id"`
		Hidden string `db:"-"`
	}

	m := sqlxutils.FieldToColMap(Row{})
	assert.NotContains(t, m, "Hidden")
}

