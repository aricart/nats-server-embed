package main

import (
	"github.com/aricart/nats-server-embed/nse"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var server *nse.NatsServer

// handle signals so we can orderly shutdown
func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		for sig := range c {
			switch sig {
			case syscall.SIGINT:
				server.Server.Shutdown()
				os.Exit(0)
			}
		}
	}()
}

func main() {
	handleSignals()
	server = nse.Start()
	if server != nil {
		runtime.Goexit()
	}
}
