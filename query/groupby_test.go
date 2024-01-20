package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_GroupBy(t *testing.T) {
	t.Run("should-create-group-by-param", func(t *testing.T) {
		g := query.GroupBy("a")
		assert.Equal(t, query.GroupByParam{Names: []string{"a"}}, g)
	})

	t.Run("should-create-group-by-param-with-option", func(t *testing.T) {
		a := query.GroupBy("a")
		b := a.WithOption("option")

		assert.NotEqual(t, a, b)

		assert.Equal(t, query.GroupByParam{
			Names:  []string{"a"},
			Option: "option",
		}, b)
	})

	t.Run("should-create-group-by-param-with-having", func(t *testing.T) {
		a := query.GroupBy("a")
		b := a.WithHaving(
			query.Filter("a", 1),
		)

		assert.NotEqual(t, a, b)

		assert.Equal(t, query.GroupByParam{
			Names:  []string{"a"},
			Having: []query.FilterParam{query.Filter("a", 1)},
		}, b)
	})
}
