package sqlxutils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sqlxutils "github.com/infevocorp/goflexstore/sqlx/utils"
)

type UserDTO struct{ ID int64 }
type OrderItemDTO struct{ ID int64 }
type HTTPRequestDTO struct{ ID int64 }
type PlainStruct struct{ ID int64 }

type customTableDTO struct{ ID int64 }

func (customTableDTO) TableName() string { return "my_custom_table" }

func TestTableName(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"strips DTO suffix and pluralises", UserDTO{}, "users"},
		{"snake_case multi-word", OrderItemDTO{}, "order_items"},
		{"snake_case acronym", HTTPRequestDTO{}, "http_requests"},
		{"no DTO suffix falls back to type name", PlainStruct{}, "plain_structs"},
		{"TableNamer interface overrides inference", customTableDTO{}, "my_custom_table"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sqlxutils.TableName(tt.input))
		})
	}
}
