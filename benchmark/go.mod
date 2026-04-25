module github.com/infevocorp/goflexstore/benchmark

go 1.21.6

replace (
	github.com/infevocorp/goflexstore => ../
	github.com/infevocorp/goflexstore/gorm => ../gorm
	github.com/infevocorp/goflexstore/sqlx => ../sqlx
)

require (
	github.com/glebarez/sqlite v1.11.0
	github.com/infevocorp/goflexstore v1.0.11
	github.com/infevocorp/goflexstore/gorm v0.0.0-00010101000000-000000000000
	github.com/infevocorp/goflexstore/sqlx v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.4.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gorm.io/hints v1.1.2 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)
