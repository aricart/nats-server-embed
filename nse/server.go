package nse

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nats-io/gnatsd/server"
	"net"
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

func Start(args []string) (*NatsServer, error) {
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet("nats-server", flag.ExitOnError)
	fs.Usage = usage
	opts, err := server.ConfigureOptions(fs, args,
		server.PrintServerAndExit,
		fs.Usage,
		server.PrintTLSHelpAndDie)

	if err != nil {
		server.PrintAndDie(err.Error() + "\n" + usageString)
	}

	s := server.New(opts)
	s.ConfigureLogger()

	go s.Start()

	if !s.ReadyForConnections(5 * time.Second) {
		return nil, errors.New("unable to start embedded server")
	}

	return &NatsServer{Server: s}, nil

}
