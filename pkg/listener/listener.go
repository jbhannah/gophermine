package listener

import (
	"context"
	"fmt"
	"net"
	"sync"

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
	stopped chan struct{}
	wg      *sync.WaitGroup
}

// NewListener creates a new listener at the given address.
func NewListener(ctx context.Context, handler Handler, addr string) (*Listener, error) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Could not listen on %s for %s: %v", addr, handler.Name(), err)
	}

	listener := &Listener{
		Handler:  handler,
		Listener: listen,
		stopped:  make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}

	listener.Runner = runner.NewRunner(ctx, listener)
	return listener, nil
}

// Setup starts the connection listening loop.
func (listener *Listener) Setup() {
	defer log.Infof("Listening on %s for %s", listener.Addr(), listener.Name())
	go listener.listen()
}

// Run restarts the listener if it stops unexpectedly.
func (listener *Listener) Run() {
	for {
		select {
		case <-listener.Done():
			return
		case <-listener.stopped:
			log.Warnf("Restarting listener for %s", listener.Name())
			go listener.listen()
		}
	}
}

// Cleanup closes the listener.
func (listener *Listener) Cleanup() {
	defer log.Debugf("Stopped listening on %s for %s", listener.Addr(), listener.Name())
	log.Debugf("Stopping listener for %s", listener.Name())

	listener.wg.Wait()
	listener.Close()
	<-listener.stopped
}

func (listener *Listener) handle(conn *net.TCPConn) {
	closed := make(chan struct{})
	defer close(closed)

	go func() {
		defer log.Debugf("Closed connection for %s from %s", listener.Name(), conn.RemoteAddr())
		defer listener.wg.Done()

		select {
		case <-listener.Done():
			if err := conn.CloseRead(); err != nil {
				log.Warnf("Unable to close connection for %s from %s nicely", listener.Name(), conn.RemoteAddr())
			}
		case <-closed:
		}
	}()

	listener.wg.Add(1)
	listener.HandleConn(conn)
}

func (listener *Listener) listen() {
	defer close(listener.stopped)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if listener.Err() == nil {
				log.Errorf("Error accepting connection for %s: %s", listener.Name(), err)
			}

			return
		}

		log.Infof("Accepted connection for %s from %s", listener.Name(), conn.RemoteAddr())
		go listener.handle(conn.(*net.TCPConn))
	}
}
