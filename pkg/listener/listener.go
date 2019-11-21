package listener

import (
	"context"
	"net"
	"os"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

// Handler defines the interface for handlers of incoming network connections.
type Handler interface {
	HandleConn(net.Conn)
}

// Listener performs non-blocking handling of incoming network connections.
type Listener struct {
	Handler
	net.Listener
	*runner.Runner
	restart chan struct{}
}

// NewListener creates a new listener at the given address.
func NewListener(ctx context.Context, handler Handler, addr string) *Listener {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Debug("Could not listen on", addr)
		os.Exit(1)
	}

	listener := &Listener{
		Handler:  handler,
		Listener: listen,
		restart:  make(chan struct{}),
	}

	listener.Runner = runner.NewRunner(ctx, listener)
	return listener
}

// Setup starts the connection listening loop.
func (listener *Listener) Setup() {
	go listener.listen()
}

// Run restarts the listener if it stops unexpectedly.
func (listener *Listener) Run() {
	for {
		select {
		case <-listener.Done():
			return
		case <-listener.restart:
			log.Warn("Restarting listener")
			go listener.listen()
		}
	}
}

// Cleanup closes the listener.
func (listener *Listener) Cleanup() {
	log.Debug("Closing listener")
	listener.Close()
}

func (listener *Listener) listen() {
	defer close(listener.restart)
	log.Debugf("Listening on %s", listener.Addr())

	for {
		log.Debug("Waiting for connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Error accepting connection: %v", err)
			break
		}

		log.Debugf("Got connection from %s", conn.RemoteAddr())
		go listener.HandleConn(conn)
	}

	log.Debug("Listener stopped")
}
