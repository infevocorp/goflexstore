package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infevocorp/goflexstore/query"
)

func Test_OrderBy(t *testing.T) {
	t.Run("param-type-should-be-orderby", func(t *testing.T) {
		assert.Equal(t, query.TypeOrderBy, query.OrderByParam{}.ParamType())
	})

	t.Run("should-create-order-by-param", func(t *testing.T) {
		o := query.OrderBy("Name", false)

		assert.Equal(t, query.OrderByParam{
			Name: "Name",
			Desc: false,
		}, o)
	})
}
