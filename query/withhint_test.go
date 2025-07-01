package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infevocorp/goflexstore/query"
)

func Test_WithHint(t *testing.T) {
	t.Run("param-type-should-be-withhint", func(t *testing.T) {
		assert.Equal(t, query.TypeWithHint, query.WithHintParam{}.ParamType())
	})

	t.Run("should-create-withhint-param", func(t *testing.T) {
		p := query.WithHint("INL_HASH_JOIN(users)")

		assert.Equal(t, query.WithHintParam{
			Hint: "INL_HASH_JOIN(users)",
		}, p)
	})
}
