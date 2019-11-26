package server

import (
	"context"
	"sync"
	"time"

	"github.com/jbhannah/gophermine/pkg/mc"

	"github.com/jbhannah/gophermine/pkg/runner"
)

// TickDuration is the length of a single world tick (50ms).
const TickDuration = 50 * time.Millisecond

// Server defines the overarching server, managing the world instance, network
// listeners, and communication between them all.
type Server struct {
	*runner.Runner
	console *Console
	mc      *MCServer
	rcon    *RCONServer
	ticker  *time.Ticker
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context) (*Server, error) {
	server := &Server{
		ticker: time.NewTicker(TickDuration),
	}

	server.Runner = runner.NewRunner(ctx, server)
	server.console, _ = NewConsole(server.Context)

	mcServer, err := NewMCServer(server.Context, mc.Properties().ServerAddr())
	if err != nil {
		return nil, err
	}

	server.mc = mcServer

	rconAddr := mc.Properties().RCONAddr()
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
	go server.console.Start()
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

	wg := &sync.WaitGroup{}

	if server.rcon != nil {
		wg.Add(3)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			<-server.rcon.Stopped()
		}(wg)
	} else {
		wg.Add(2)
	}

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-server.mc.Stopped()
	}(wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-server.console.Stopped()
	}(wg)

	wg.Wait()
}
