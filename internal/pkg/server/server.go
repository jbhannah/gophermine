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
	mc     *MCServer
	rcon   *RCONServer
	ticker *time.Ticker
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, mcAddr string, rconAddr string) (*Server, error) {
	server := &Server{
		ticker: time.NewTicker(TickDuration),
	}

	server.Runner = runner.NewRunner(ctx, server)

	mc, err := NewMCServer(server.Context, mcAddr)
	if err != nil {
		return nil, err
	}

	server.mc = mc

	if rconAddr != "" {
		rcon, err := NewRCONServer(server.Context, rconAddr)
		if err != nil {
			return nil, err
		}

		server.rcon = rcon
	}

	return server, nil
}

// Name returns the name of the server.
func (server *Server) Name() string {
	return "Gophermine"
}

// Setup starts the server's network listeners.
func (server *Server) Setup() {
	go server.mc.Start()

	if server.rcon != nil {
		go server.rcon.Start()
	}
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

	if server.rcon != nil {
		<-server.rcon.Stopped()
	}
}
