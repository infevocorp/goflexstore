package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_OrderBy(t *testing.T) {
	g := query.OrderBy("a", true)
	assert.Equal(t, query.OrderByParam{Name: "a", Desc: true}, g)
}
