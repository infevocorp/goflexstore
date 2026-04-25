package sqlxquery

import "github.com/jmoiron/sqlx"

// Dialect identifies the target SQL dialect used for placeholder rebinding
// and dialect-specific SQL in Upsert.
type Dialect string

const (
	DialectMySQL    Dialect = "mysql"
	DialectPostgres Dialect = "postgres"
	DialectSQLite   Dialect = "sqlite"
)

// BindType returns the sqlx bind type constant for the dialect.
func (d Dialect) BindType() int {
	if d == DialectPostgres {
		return sqlx.DOLLAR
	}
	return sqlx.QUESTION
}

// Rebind rewrites a query that uses `?` placeholders into the dialect-specific
// form (e.g. `$1`, `$2` for Postgres).
func Rebind(dialect Dialect, query string) string {
	return sqlx.Rebind(dialect.BindType(), query)
}
