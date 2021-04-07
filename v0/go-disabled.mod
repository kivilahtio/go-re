## For some strange reason Semantic Import Versioning doesn't work with versions v0 and v1, only v2 onwards.
## Skipping submodules seems to work tho
module github.com/kivilahtio/go-re/v1

go 1.15

require (
	github.com/smartystreets/goconvey v1.6.4 // test
)
