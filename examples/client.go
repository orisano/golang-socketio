package main

import (
	"log"
	"runtime"

	"github.com/orisano/golang-socketio"
	"github.com/orisano/golang-socketio/transport"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetURL("localhost:3811", false),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
	}

	c.Methods.On("/message", func(h *gosocketio.Channel, args string) interface{} {
		return nil
	})
	c.Methods.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel, args string) interface{} {
		return nil
	})
	c.Methods.On(gosocketio.OnConnection, func(h *gosocketio.Channel, arg string) interface{} {
		return nil
	})
}
