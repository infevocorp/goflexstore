package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goflexstore/query"
)

func Test_WithLock(t *testing.T) {
	t.Run("param-type-should-be-withlock", func(t *testing.T) {
		assert.Equal(t, query.TypeWithLock, query.WithLockParam{}.ParamType())
	})

	t.Run("should-create-withlock-param", func(t *testing.T) {
		p := query.WithLock(query.LockTypeForUpdate)

		assert.Equal(t, query.WithLockParam{
			LockType: query.LockTypeForUpdate,
		}, p)
	})
}
