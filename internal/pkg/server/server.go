package server

import (
	"context"
	"time"
)

// TickDuration is the length of a single world tick (50ms).
const TickDuration = 50 * time.Millisecond

// Server defines the overarching server, managing the world instance, network
// listeners, and communication between them all.
type Server struct {
	context.Context
	*time.Ticker
	Stop context.CancelFunc
	done context.CancelFunc
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, done context.CancelFunc) *Server {
	ctx, cancel := context.WithCancel(ctx)

	return &Server{
		Context: ctx,
		Ticker:  time.NewTicker(TickDuration),
		Stop:    cancel,
		done:    done,
	}
}

// Start starts the server.
func (server *Server) Start() {
	println("Starting server")
	go server.run()

	println("Started server")
}

func (server *Server) run() {
	defer server.cleanup()

	for {
		select {
		case <-server.Done():
			println("Stopping server")
			return
		}
	}
}

func (server *Server) cleanup() {
	defer server.done()
	println("Stopped server")
}
