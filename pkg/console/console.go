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

type Console struct {
	*runner.Runner
	*bufio.Scanner
	commands   chan *mc.Command
	ctxStarted chan struct{}
	lines      chan string
	name       string
	input      io.Reader
}

func NewConsole(ctx context.Context, name string, input io.Reader) (*Console, error) {
	console := &Console{
		Scanner:    bufio.NewScanner(input),
		commands:   ctx.Value(mc.ServerCommands).(chan *mc.Command),
		ctxStarted: ctx.Value(runner.RunnableStarted).(chan struct{}),
		lines:      make(chan string),
		name:       name,
		input:      input,
	}

	console.Runner = runner.NewRunner(ctx, console)
	return console, nil
}

func (console *Console) Name() string {
	return console.name
}

func (console *Console) Setup() {
	go console.scan()
}

func (console *Console) Run() {
	<-console.ctxStarted
	log.Debug("Accepting console commands")

	for {
		select {
		case <-console.Done():
			return
		case text := <-console.lines:
			console.commands <- mc.NewCommand(console, strings.Split(text, " ")...)
		}
	}
}

func (console *Console) Cleanup() {
	defer close(console.lines)
}

func (console *Console) scan() {
	for console.Scan() {
		console.lines <- console.Text()
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}
