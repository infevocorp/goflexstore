package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_GroupBy(t *testing.T) {
	g := query.GroupBy("a")
	assert.Equal(t, query.GroupByParam{Name: "a"}, g)
}

func Test_GroupByOption(t *testing.T) {
	g := query.GroupByOption("b")
	assert.Equal(t, query.GroupByOptionParam{Option: "b"}, g)
}
