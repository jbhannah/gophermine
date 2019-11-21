package mc

import (
	"context"
	"io"
	"net"
	"os"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

// Listener performs non-blocking handling of incoming network connections.
type Listener struct {
	net.Listener
	*runner.Runner
	listening chan struct{}
}

// NewListener creates a new listener at the given address.
func NewListener(ctx context.Context, addr string) *Listener {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Debug("Could not listen on", addr)
		os.Exit(1)
	}

	listener := &Listener{
		Listener:  listen,
		listening: make(chan struct{}),
	}

	listener.Runner = runner.NewRunner(ctx, listener)
	return listener
}

// Setup starts the connection listening loop.
func (listener *Listener) Setup() {
	log.Debugf("Listening on %s", listener.Addr())
	go listener.listen()
}

// Run handles incoming connections and restarts the listener if it stops
// unexpectedly.
func (listener *Listener) Run() {
	select {
	case <-listener.listening:
		log.Warn("Restarting listener")
		go listener.listen()
	default:
		return
	}
}

// Cleanup closes the listener.
func (listener *Listener) Cleanup() {
	log.Debug("Closing listener")
	listener.Close()
}

func (listener *Listener) listen() {
	defer close(listener.listening)
	log.Debug("Listening")

	for {
		log.Debug("Waiting for connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Error accepting connection: %v", err)
			break
		}

		log.Debugf("Accepted connection from %s", conn.RemoteAddr())
		go io.Copy(conn, conn)
	}

	log.Debug("Listener stopped")
}
