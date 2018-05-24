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
	var err error

	inlineArgs := -1
	for i, a := range os.Args {
		if a == "--" {
			inlineArgs = i + 1
			break
		}
	}

	var args []string
	if inlineArgs > -1 {
		args = os.Args[inlineArgs+1:]
	}

	server, err = nse.Start(args)
	if err != nil {
		panic(err)
	}
	if server != nil {
		runtime.Goexit()
	}
}
