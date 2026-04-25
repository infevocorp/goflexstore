package sqlxutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"UserName", "user_name"},
		{"OrderItem", "order_item"},
		{"HTTPRequest", "http_request"},
		{"MyHTTPClient", "my_http_client"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"already_snake", "already_snake"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toSnakeCase(tt.input))
		})
	}
}
