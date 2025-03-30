module github.com/lexcao/genapi/pkg/clients/resty

go 1.21

replace github.com/lexcao/genapi => ../../../

require (
	github.com/go-resty/resty/v2 v2.16.5
	github.com/lexcao/genapi v0.0.1
)

require (
	github.com/jarcoal/httpmock v1.3.1 // indirect
	golang.org/x/net v0.33.0 // indirect
)
