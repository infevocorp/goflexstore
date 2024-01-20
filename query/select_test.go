package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_Select(t *testing.T) {
	p := query.Select("id", "name")

	assert.Equal(t,
		query.SelectParam{
			Names: []string{"id", "name"},
		},
		p,
	)
}
