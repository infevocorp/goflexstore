package gormquery_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormquery "github.com/jkaveri/goflexstore/gorm/query"
	gormutils "github.com/jkaveri/goflexstore/gorm/utils"
	"github.com/jkaveri/goflexstore/query"
)

type User struct {
	ID        int    `gorm:"column:id;primary_key;auto_increment"`
	Name      string `gorm:"column:name"`
	Age       int    `gorm:"column:age"`
	RefererID int    `gorm:"column:referer_id"`
	Referer   *User  `gorm:"foreignKey:RefererID"`
}

func Test_Builder_Build(t *testing.T) {
	type deps struct {
		sql sqlmock.Sqlmock
	}

	type args struct {
		params query.Params
	}

	type expects struct {
		err   bool
		users []User
	}

	tests := []struct {
		name    string
		args    args
		expects expects
		mock    func(d deps)
	}{
		{
			name: "filter",
			args: args{
				params: query.NewParams(
					query.Filter("name", "john"),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE name = ?")).
					WithArgs("john").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "filter-name-and-age",
			args: args{
				params: query.NewParams(
					query.Filter("name", "john"),
					query.Filter("age", 20),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE name = ? AND age = ?")).
					WithArgs("john", 20).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "filter-name-or",
			args: args{
				params: query.NewParams(
					query.OR(query.Filter("name", "john"), query.Filter("name", "jenny")),
					query.Filter("age", 20),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (name = ? OR name = ?) AND age = ?")).
					WithArgs("john", "jenny", 20).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "paginate",
			args: args{
				params: query.NewParams(
					query.Paginate(1, 10),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` LIMIT 10 OFFSET 1")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "order-by",
			args: args{
				params: query.NewParams(
					query.OrderBy("Name", true),
					query.OrderBy("ID", false),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` ORDER BY `name` DESC,`id`")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "group-by",
			args: args{
				params: query.NewParams(
					query.GroupBy("Name"),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` GROUP BY `name`")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(1, "john", 20))
			},
		},

		{
			name: "select",
			args: args{
				params: query.NewParams(
					query.Select("Name", "Age"),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   0,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT `name`,`age` FROM `users`")).
					WillReturnRows(sqlmock.NewRows([]string{"name", "age"}).
						AddRow("john", 20))
			},
		},

		{
			name: "preload",
			args: args{
				params: query.NewParams(
					query.Filter("RefererID", 0).WithOP(query.NEQ),
					query.Preload("Referer"),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:        1,
						Name:      "john",
						Age:       20,
						RefererID: 2,
						Referer: &User{
							ID:   2,
							Name: "jenny",
							Age:  20,
						},
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE referer_id <> ?")).
					WithArgs(0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "referer_id"}).
						AddRow(1, "john", 20, 2))

				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(2, "jenny", 20))
			},
		},

		{
			name: "preload-with-filter",
			args: args{
				params: query.NewParams(
					query.Filter("RefererID", 0).WithOP(query.NEQ),
					query.Preload("Referer",
						query.Filter("Name", "jenny"),
						query.Filter("Age", 20),
					),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:        1,
						Name:      "john",
						Age:       20,
						RefererID: 2,
						Referer: &User{
							ID:   2,
							Name: "jenny",
							Age:  20,
						},
					},
				},
			},
			mock: func(d deps) {
				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE referer_id <> ?")).
					WithArgs(0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age", "referer_id"}).
						AddRow(1, "john", 20, 2))

				d.sql.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE name = ? AND age = ? AND `users`.`id` = ?")).
					WithArgs("jenny", 20, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).
						AddRow(2, "jenny", 20))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := newTestDB(t)

			d := deps{
				sql: sqlMock,
			}

			tt.mock(d)

			builder := gormquery.NewBuilder(
				gormquery.WithFieldToColMap(gormutils.FieldToColMap(User{})),
			)
			scopes := builder.Build(tt.args.params)

			var users []User
			err := db.Scopes(scopes...).Find(&users).Error

			require.Equal(t, tt.expects.err, err != nil, "unepxected error: %v", err)
			require.Equal(t, tt.expects.users, users)
		})
	}
}

func Test_ScopeBuilder_CustomFilter(t *testing.T) {
	type deps struct {
		sql sqlmock.Sqlmock
	}

	type args struct {
		customFilters map[string]gormquery.ScopeBuilderFunc
		params        query.Params
	}

	type expects struct {
		err   bool
		users []User
	}

	tests := []struct {
		name    string
		args    args
		expects expects
		mock    func(d deps)
	}{
		{
			name: "custom-filter-should-be-called",
			args: args{
				customFilters: map[string]gormquery.ScopeBuilderFunc{
					"name": func(param query.Param) gormquery.ScopeFunc {
						return func(tx *gorm.DB) *gorm.DB {
							p := param.(query.FilterParam)

							return tx.Where("`first_name` = ?", p.Value)
						}
					},
				},
				params: query.NewParams(
					query.Filter("name", "john"),
				),
			},
			expects: expects{
				err: false,
				users: []User{
					{
						ID:   1,
						Name: "john",
						Age:  20,
					},
				},
			},
			mock: func(d deps) {
				d.sql.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `first_name` = ?")).
					WithArgs("john").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "john", 20))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, sqlMock := newTestDB(t)

			d := deps{
				sql: sqlMock,
			}

			tt.mock(d)

			builder := gormquery.NewBuilder(
				gormquery.WithCustomFilters(tt.args.customFilters),
			)
			scopes := builder.Build(tt.args.params)

			var users []User
			err := db.Scopes(scopes...).Find(&users).Error

			assert.Equal(t, tt.expects.err, err != nil, "unepxected error: %v", err)
			assert.Equal(t, tt.expects.users, users)
		})
	}
}

func newTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlMock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.23"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})

	t.Cleanup(func() {
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	return gormDB, sqlMock
}
