package console

import (
	"bufio"
	"context"
	"io"

	"github.com/jbhannah/gophermine/pkg/runner"

	log "github.com/sirupsen/logrus"
)

type Console struct {
	*runner.Runner
	*bufio.Scanner
	ctxStarted chan struct{}
	input      chan string
	name       string
	stream     io.Reader
}

func NewConsole(ctx context.Context, name string, stream io.Reader) (*Console, error) {
	console := &Console{
		Scanner:    bufio.NewScanner(stream),
		ctxStarted: ctx.Value(runner.RunnableStarted).(chan struct{}),
		input:      make(chan string),
		name:       name,
		stream:     stream,
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
		case text := <-console.input:
			log.Info(text)
		}
	}
}

func (console *Console) scan() {
	for console.Scan() {
		console.input <- console.Text()
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}

func (console *Console) Cleanup() {
	defer close(console.input)
}
