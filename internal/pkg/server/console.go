package server

import (
	"bufio"
	"context"
	"os"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

type Console struct {
	*runner.Runner
	*bufio.Scanner
	serverStarted chan struct{}
}

func NewConsole(ctx context.Context) (*Console, error) {
	console := &Console{
		Scanner:       bufio.NewScanner(os.Stdin),
		serverStarted: ctx.Value(runner.RunnableStarted).(chan struct{}),
	}

	console.Runner = runner.NewRunner(ctx, console)
	return console, nil
}

func (console *Console) Name() string {
	return "Console"
}

func (console *Console) Setup() {}

func (console *Console) Run() {
	<-console.serverStarted
	log.Debug("Accepting console commands")
	go console.Scan()
	<-console.Done()
}

func (console *Console) Cleanup() {
	os.Stdin.Close()
}

func (console *Console) Scan() {
	for console.Scanner.Scan() {
		log.Infof("%s", console.Text())
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}
