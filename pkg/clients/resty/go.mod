module github.com/lexcao/genapi/pkg/clients/resty

go 1.21
toolchain go1.24.1

replace github.com/lexcao/genapi => ../../../

require (
	github.com/go-resty/resty/v2 v2.16.5
	github.com/lexcao/genapi v0.0.1
)

require golang.org/x/net v0.36.0 // indirect
