package server

import (
	"context"
	"time"

	"github.com/jbhannah/gophermine/pkg/runner"
)

// TickDuration is the length of a single world tick (50ms).
const TickDuration = 50 * time.Millisecond

// Server defines the overarching server, managing the world instance, network
// listeners, and communication between them all.
type Server struct {
	*runner.Runner
	rcon   *RCONServer
	ticker *time.Ticker
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, addr string) *Server {
	server := &Server{
		ticker: time.NewTicker(TickDuration),
	}

	server.Runner = runner.NewRunner(ctx, server)
	server.rcon = NewRCONServer(server.Context, addr)
	return server
}

// Setup starts the server's network listeners.
func (server *Server) Setup() {
	server.rcon.Start()
}

// Run handles incoming commands to the server.
func (server *Server) Run() {
	for {
		select {
		case <-server.Done():
			return
		case <-server.ticker.C:
		}
	}
}

// Cleanup stops the server's ticker and network listeners.
func (server *Server) Cleanup() {
	server.ticker.Stop()
	<-server.rcon.Stopped()
}
