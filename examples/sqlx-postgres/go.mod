module github.com/infevocorp/goflexstore/examples/sqlx-postgres

go 1.21.6

replace (
	github.com/infevocorp/goflexstore => ../../
	github.com/infevocorp/goflexstore/sqlx => ../../sqlx
)

require (
	github.com/infevocorp/goflexstore v1.0.11
	github.com/infevocorp/goflexstore/sqlx v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
