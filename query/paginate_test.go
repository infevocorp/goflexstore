package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_Paginate(t *testing.T) {
	t.Run("param-type-should-be-paginate", func(t *testing.T) {
		assert.Equal(t, query.TypePaginate, query.PaginateParam{}.ParamType())
	})

	t.Run("should-create-paginate-param", func(t *testing.T) {
		p := query.Paginate(1, 2)

		assert.Equal(t, query.PaginateParam{
			Offset: 1,
			Limit:  2,
		}, p)
	})
}
