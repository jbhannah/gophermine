package mc

import (
	"context"
	"io"
	"net"
	"os"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

type Listener struct {
	net.Listener
	*runner.Runner
}

func NewListener(ctx context.Context, addr string) *Listener {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Debug("Could not listen on", addr)
		os.Exit(1)
	}

	listener := &Listener{
		Listener: listen,
	}

	listener.Runner = runner.NewRunner(ctx, listener)
	return listener
}

func (listener *Listener) Setup() {
	log.Debug("Listening")
	go listener.listen()
}

func (listener *Listener) Cleanup() {
	log.Debug("Closing listener")
	listener.Close()
}

func (listener *Listener) listen() {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Debug("Error:", err)
			return
		}

		go io.Copy(conn, conn)
	}
}
