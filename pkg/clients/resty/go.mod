module github.com/lexcao/genapi/pkg/clients/resty

go 1.23.0

toolchain go1.24.5

replace github.com/lexcao/genapi => ../../../

require (
	github.com/go-resty/resty/v2 v2.16.5
	github.com/lexcao/genapi v0.0.1
)

require golang.org/x/net v0.38.0 // indirect
