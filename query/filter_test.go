package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infevocorp/goflexstore/query"
)

func Test_WithOperator(t *testing.T) {
	ageFilter := query.Filter("age", 10)
	assert.Equal(t, query.FilterParam{
		Name:     "age",
		Operator: query.EQ,
		Value:    10,
	}, ageFilter)

	t.Run("EQ", func(t *testing.T) {
		param := ageFilter.WithOP(query.GT)

		assert.Equal(t, query.FilterParam{
			Name:     "age",
			Operator: query.GT,
			Value:    10,
		}, param)
	})

	t.Run("LTE", func(t *testing.T) {
		param := ageFilter.WithOP(query.LTE)

		assert.Equal(t, query.FilterParam{
			Name:     "age",
			Operator: query.LTE,
			Value:    10,
		}, param)
	})
}

func Test_Filter(t *testing.T) {
	t.Run("EQ", func(t *testing.T) {
		param := query.Filter("name", "john")

		assert.Equal(t, query.FilterParam{
			Name:     "name",
			Operator: query.EQ,
			Value:    "john",
		}, param)
	})
}
