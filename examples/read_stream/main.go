package main

// This file provides an example of how to initiate and terminate a stream of
// readings using the WebSocket client. In this example, readings are printed
// out to console. After 6 seconds, the stream is terminated and the program
// will exit.

import (
	"fmt"
	"log"
	"time"

	"github.com/vapor-ware/synse-client-go/synse"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
)

func main() {
	c, err := synse.NewWebSocketClientV3(&synse.Options{
		Address: "localhost:5000",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
	if err := c.Open(); err != nil {
		log.Fatal(err)
	}

	stop := make(chan struct{}, 1)
	readings := make(chan *scheme.Read)
	defer close(readings)

	go func() {
		if err := c.ReadStream(scheme.ReadStreamOptions{}, readings, stop); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("-- Streaming Readings --")
	timeout := time.After(6 * time.Second)
	for {
		select {
		case r := <-readings:
			fmt.Printf("• %+v\n", r)

		case <-timeout:
			fmt.Println("-- terminating stream --")
			close(stop)
			return
		}
	}
}
