package server

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/jbhannah/gophermine/pkg/console"
	"github.com/jbhannah/gophermine/pkg/mc"
	"github.com/jbhannah/gophermine/pkg/runner"

	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
)

// TickDuration is the length of a single world tick (50ms).
const TickDuration = 50 * time.Millisecond

// Server defines the overarching server, managing the world instance, network
// listeners, and communication between them all.
type Server struct {
	*runner.Runner
	commands <-chan *mc.Command
	console  *console.Console
	mc       *MCServer
	rcon     *RCONServer
	ticker   *time.Ticker
}

// NewServer instantiates a new server.
func NewServer(ctx context.Context) (*Server, error) {
	cmds := make(chan *mc.Command)

	server := &Server{
		commands: cmds,
		ticker:   time.NewTicker(TickDuration),
	}

	ctx = context.WithValue(ctx, mc.ServerCommands, cmds)
	server.Runner = runner.NewRunner(ctx, server)

	if isatty.IsTerminal(os.Stdin.Fd()) {
		cons, err := console.NewConsole(server.Context, "Console", os.Stdin)
		if err != nil {
			return nil, err
		}

		server.console = cons
	}

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
	wg := &sync.WaitGroup{}
	wg.Add(1)

	if server.rcon != nil {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			<-server.rcon.Start()
		}(wg)
	}

	if server.console != nil {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			<-server.console.Start()
		}(wg)
	}

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-server.mc.Start()
	}(wg)

	wg.Wait()
}

// Run handles incoming commands to the server.
func (server *Server) Run() {
	for {
		select {
		case <-server.Done():
			return
		case cmd := <-server.commands:
			go server.handleCommand(cmd)
		case <-server.ticker.C:
		}
	}
}

// Cleanup stops the server's ticker and network listeners.
func (server *Server) Cleanup() {
	server.ticker.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	if server.rcon != nil {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			<-server.rcon.Stopped()
		}(wg)
	}

	if server.console != nil {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			<-server.console.Stopped()
		}(wg)
	}

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-server.mc.Stopped()
	}(wg)

	wg.Wait()
}

func (server *Server) handleCommand(cmd *mc.Command) {
	log.Debugf("Command received: %s", cmd)

	switch cmd.CommandType {
	case mc.StopCommand:
		server.Stop()
	}
}
