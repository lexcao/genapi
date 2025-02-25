package main

import (
	"fmt"

	"github.com/lexcao/genapi"
	"github.com/lexcao/genapi/examples/httpbin/api"
)

func main() {
	client := genapi.New[api.HttpBin]()
	resp, err := client.Get("World")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Args)
}
