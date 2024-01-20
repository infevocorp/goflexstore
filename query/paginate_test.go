package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_Paginate(t *testing.T) {
	p := query.Paginate(10, 20)
	assert.Equal(t, query.PaginateParam{
		Offset: 10,
		Limit:  20,
	}, p)
}
