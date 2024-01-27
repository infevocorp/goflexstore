module github.com/jkaveri/goflexstore/examples/cms

go 1.21.6

replace (
	github.com/jkaveri/goflexstore => ../..
	github.com/jkaveri/goflexstore/gorm => ../../gorm
)

require (
	github.com/jkaveri/goflexstore v1.0.2
	github.com/jkaveri/goflexstore/gorm v1.0.2
)

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/labstack/echo/v4 v4.11.4
	github.com/pkg/errors v0.9.1 // indirect
	gorm.io/gorm v1.25.5 // indirect
)
