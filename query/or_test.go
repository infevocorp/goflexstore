package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infevocorp/goflexstore/query"
)

func Test_OR(t *testing.T) {
	t.Run("param-type-should-be-or", func(t *testing.T) {
		assert.Equal(t, query.TypeOR, query.ORParam{}.ParamType())
	})

	t.Run("should-create-or-param", func(t *testing.T) {
		o := query.OR(
			query.Filter("id", 1),
			query.Filter("id", 2),
		)

		assert.Equal(t, query.ORParam{
			Params: []query.FilterParam{
				query.Filter("id", 1),
				query.Filter("id", 2),
			},
		}, o)
	})

	t.Run("should-panic-if-param-is-not-filter", func(t *testing.T) {
		assert.Panics(t, func() {
			query.OR(
				query.Filter("id", 1),
				query.GroupBy("id"),
			)
		})
	})
}
