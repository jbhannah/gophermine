package server

import (
	"context"
	"time"

	"github.com/jbhannah/gophermine/internal/pkg/mc"

	"github.com/jbhannah/gophermine/pkg/runner"
)

// TickDuration is the length of a single world tick (50ms).
const TickDuration = 50 * time.Millisecond

// Server defines the overarching server, managing the world instance, network
// listeners, and communication between them all.
type Server struct {
	*runner.Runner
	listener *mc.Listener
	ticker   *time.Ticker
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context, addr string) *Server {
	server := &Server{
		listener: mc.NewListener(ctx, addr),
		ticker:   time.NewTicker(TickDuration),
	}

	return runner.NewRunner(ctx, server).(*Server)
}

func (server *Server) WithRunner(runner *runner.Runner) runner.Runnable {
	server.Runner = runner
	return server
}

func (server *Server) Setup() {
	server.listener.Start()
}

func (server *Server) Cleanup() {
	server.ticker.Stop()
	<-server.listener.Stop()
}
