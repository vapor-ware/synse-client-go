package main

// http.go shows a simple example of a http client's usage.

import (
	"fmt"

	"github.com/vapor-ware/synse-client-go/synse"
)

func main() {
	config := &synse.Options{
		Server: synse.ServerOptions{
			Address: "localhost:5000",
		},
	}

	client, err := synse.NewHTTPClient(config)
	fmt.Printf("%+v, %+v\n", client, err)

	r1, err := client.Status()
	fmt.Printf("%+v, %+v\n", r1, err)

	r2, err := client.Version()
	fmt.Printf("%+v, %+v\n", r2, err)

	r3, err := client.Config()
	fmt.Printf("%+v, %+v\n", r3, err)
}
