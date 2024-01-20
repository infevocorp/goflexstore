package gormstore_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	gormopscope "github.com/jkaveri/goflexstore/gorm/opscope"
	gormstore "github.com/jkaveri/goflexstore/gorm/store"
	"github.com/jkaveri/goflexstore/query"
)

func Test_Store_Get(t *testing.T) {
	type args struct {
		ctx    context.Context
		params []query.Param
	}

	type expecteds struct {
		err  bool
		user User
	}

	type deps struct {
		sqlMock sqlmock.Sqlmock
	}

	tests := []struct {
		name string
		args args
		mock func(deps)
		want expecteds
	}{
		{
			name: "get-by-id",
			args: args{},
			mock: func(deps) {
			},
			want: expecteds{},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := newTestDB(t)

			d := deps{
				sqlMock: sqlMock,
			}
			tt.mock(d)

			s := gormstore.Store[User, UserDTO, int]{
				OpScope: gormopscope.NewTransactionScope(
					"test",
					db, &sql.TxOptions{
						Isolation: sql.LevelDefault,
						ReadOnly:  false,
					},
				),
			}

			got, err := s.Get(tt.args.ctx, tt.args.params...)
			assert.Equal(t, tt.want.err, err != nil)
			assert.Equal(t, tt.want.user, got)
		})
	}
}
