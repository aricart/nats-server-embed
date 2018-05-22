package nse

import (
	"flag"
	"fmt"
	"github.com/nats-io/gnatsd/server"
	"net"
	"os"
	"time"
)

type NatsServer struct {
	Server *server.Server
}

var usageString = "embedded NATS server options can be supplied by following a '--' argument with any gnatsd supported flag."

func usage() {
	fmt.Println(usageString)
}

func (nse *NatsServer) getServerPort() int {
	return nse.Server.Addr().(*net.TCPAddr).Port
}

func Start() *NatsServer {
	inlineArgs := -1
	for i, a := range os.Args {
		if a == "--" {
			inlineArgs = i + 1
			break
		}
	}

	if inlineArgs > 0 {
		// Create a FlagSet and sets the usage
		fs := flag.NewFlagSet("nats-server", flag.ExitOnError)
		fs.Usage = usage
		opts, err := server.ConfigureOptions(fs, os.Args[inlineArgs:],
			server.PrintServerAndExit,
			fs.Usage,
			server.PrintTLSHelpAndDie)

		if err != nil {
			server.PrintAndDie(err.Error() + "\n" + usageString)
		}

		s := server.New(opts)
		s.ConfigureLogger()

		fmt.Println("Starting gnatsd")

		go s.Start()

		if !s.ReadyForConnections(5 * time.Second) {
			panic("unable to start embedded server")
		}

		return &NatsServer{Server: s}
	} else {
		fmt.Println("No gnatsd configuration/arguments provided")
	}

	return nil
}
