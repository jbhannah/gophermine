package console

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/jbhannah/gophermine/pkg/mc"
	"github.com/jbhannah/gophermine/pkg/runner"

	log "github.com/sirupsen/logrus"
)

// Console is a command entry console for a running server. If the server is
// started with an attached TTY, one is instantiated to accept directly entered
// commands. If RCON is enabled for the server, one is instantiated for each
// incoming connection.
type Console struct {
	io.Reader
	io.Writer
	*runner.Runner
	*bufio.Scanner
	commands   chan *mc.Command
	ctxStarted chan struct{}
	name       string
}

// NewConsole creates a new console.
func NewConsole(ctx context.Context, name string, reader io.Reader, writer io.Writer) (*Console, error) {
	console := &Console{
		Reader:     reader,
		Writer:     writer,
		Scanner:    bufio.NewScanner(reader),
		commands:   ctx.Value(mc.ServerCommands).(chan *mc.Command),
		ctxStarted: ctx.Value(runner.RunnableStarted).(chan struct{}),
		name:       name,
	}

	console.Runner = runner.NewRunner(ctx, console)
	return console, nil
}

// Name returns the name of the console.
func (console *Console) Name() string {
	return console.name
}

// Setup begins the input scanner loop for the console.
func (console *Console) Setup() {
	go console.scan()
}

// Run waits until the surrounding context has started, then blocks until the
// console is stopped.
func (console *Console) Run() {
	<-console.ctxStarted
	log.Debugf("Accepting console commands from %s", console.Name())

	<-console.Done()
}

func (console *Console) Cleanup() {}

func (console *Console) scan() {
	for console.Scan() {
		console.commands <- mc.NewCommand(console, strings.Split(console.Text(), " ")...)
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}
