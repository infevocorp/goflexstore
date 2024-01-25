package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_Params_Get(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		filterParams := params.Get("filter")

		assert.Equal(t, []query.Param{
			query.Filter("name", "john"),
		}, filterParams)
	})

	t.Run("multiple", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
			query.Filter("age", 20),
		)

		filterParams := params.Get("filter")

		assert.Equal(t, []query.Param{
			query.Filter("name", "john"),
			query.Filter("age", 20),
		}, filterParams)
	})

	t.Run("notfound", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		filterParams := params.Get("group")

		assert.Equal(t, []query.Param{}, filterParams)
	})

	t.Run("empty", func(t *testing.T) {
		params := query.NewParams()

		filterParams := params.Get("group")

		assert.Equal(t, []query.Param{}, filterParams)
	})
}

func Test_Params_GetFilter(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		filterParam, ok := params.GetFilter("name")

		assert.True(t, ok)
		assert.Equal(t, query.Filter("name", "john"), filterParam)
	})

	t.Run("multiple", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
			query.Filter("age", 20),
		)

		filterParam, ok := params.GetFilter("age")

		assert.True(t, ok)
		assert.Equal(t, query.Filter("age", 20), filterParam)
	})

	t.Run("notfound", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		_, ok := params.GetFilter("group")

		assert.False(t, ok)
	})

	t.Run("empty", func(t *testing.T) {
		params := query.NewParams()

		_, ok := params.GetFilter("group")

		assert.False(t, ok)
	})
}

func Test_Params(t *testing.T) {
	t.Run("should-return-params", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		assert.Equal(t, []query.Param{
			query.Filter("name", "john"),
		}, params.Params())
	})
}

func Test_Params_FilterGetter(t *testing.T) {
	t.Run("should-return-filter-getter", func(t *testing.T) {
		params := query.NewParams(
			query.Filter("name", "john"),
		)

		getter := query.FilterGetter("name")

		filterParam, ok := getter(params)

		assert.True(t, ok)
		assert.Equal(t, query.Filter("name", "john"), filterParam)
	})
}
