package listener

import (
	"context"
	"net"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

// Handler defines the interface for handlers of incoming network connections.
type Handler interface {
	HandleConn(net.Conn)
	Name() string
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
		log.Fatalf("Could not listen on %s for %s", addr, handler.Name())
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
			log.Warnf("Restarting listener for %s", listener.Name())
			go listener.listen()
		}
	}
}

// Cleanup closes the listener.
func (listener *Listener) Cleanup() {
	listener.Close()
}

func (listener *Listener) handle(conn net.Conn) {
	defer log.Debugf("Closed connection for %s from %s", listener.Name(), conn.RemoteAddr())
	listener.HandleConn(conn)
}

func (listener *Listener) listen() {
	defer close(listener.restart)
	log.Infof("Listening on %s for %s", listener.Addr(), listener.Name())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Error accepting connection for %s: %v", listener.Name(), err)
			break
		}

		log.Infof("Accepted connection for %s from %s", listener.Name(), conn.RemoteAddr())
		go listener.handle(conn)
	}

	log.Debugf("Stopped listening on %s for %s", listener.Addr(), listener.Name())
}
