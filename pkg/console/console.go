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
	reader     io.Reader
	ctxStarted chan struct{}
}

func NewConsole(ctx context.Context, reader io.Reader) (*Console, error) {
	console := &Console{
		Scanner:    bufio.NewScanner(reader),
		reader:     reader,
		ctxStarted: ctx.Value(runner.RunnableStarted).(chan struct{}),
	}

	console.Runner = runner.NewRunner(ctx, console)
	return console, nil
}

func (console *Console) Name() string {
	return "Console"
}

func (console *Console) Setup() {}

func (console *Console) Run() {
	<-console.ctxStarted
	log.Debug("Accepting console commands")
	go console.Scan()
	<-console.Done()
}

func (console *Console) Cleanup() {}

func (console *Console) Scan() {
	for console.Scanner.Scan() {
		log.Infof("%s", console.Text())
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}
