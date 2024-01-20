package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_Preload_Test(t *testing.T) {
	t.Run("preload", func(t *testing.T) {
		a := query.Preload("User")
		assert.Equal(t, query.PreloadParam{
			Preload: "User",
			Params:  nil,
		}, a)
	})

	t.Run("preload-with-params", func(t *testing.T) {
		a := query.Preload("Comments", query.Filter("disabled", false))
		assert.Equal(t, query.PreloadParam{
			Preload: "Comments",
			Params: []query.Param{
				query.Filter("disabled", false),
			},
		}, a)
	})
}
