package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infevocorp/goflexstore/query"
)

func Test_Select(t *testing.T) {
	t.Run("param-type-should-be-select", func(t *testing.T) {
		assert.Equal(t, query.TypeSelect, query.SelectParam{}.ParamType())
	})

	t.Run("should-create-select-param", func(t *testing.T) {
		s := query.Select("a", "b")

		assert.Equal(t, query.SelectParam{
			Names: []string{"a", "b"},
		}, s)
	})
}
