package query_test

import (
	"testing"

	"github.com/jkaveri/goflexstore/query"
	"github.com/stretchr/testify/assert"
)

func Test_Operator_String(t *testing.T) {
	t.Run("EQ", func(t *testing.T) {
		assert.Equal(t, "EQ", query.EQ.String())
	})

	t.Run("NEQ", func(t *testing.T) {
		assert.Equal(t, "NEQ", query.NEQ.String())
	})

	t.Run("GT", func(t *testing.T) {
		assert.Equal(t, "GT", query.GT.String())
	})

	t.Run("GTE", func(t *testing.T) {
		assert.Equal(t, "GTE", query.GTE.String())
	})

	t.Run("LT", func(t *testing.T) {
		assert.Equal(t, "LT", query.LT.String())
	})

	t.Run("LTE", func(t *testing.T) {
		assert.Equal(t, "LTE", query.LTE.String())
	})

	t.Run("UNKNOWN", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN", query.Operator(100).String())
	})
}
