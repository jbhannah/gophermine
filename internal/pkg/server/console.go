package server

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/jbhannah/gophermine/pkg/runner"
	log "github.com/sirupsen/logrus"
)

type Console struct {
	*runner.Runner
	*bufio.Scanner
}

func NewConsole(ctx context.Context) (*Console, error) {
	console := &Console{
		Scanner: bufio.NewScanner(os.Stdin),
	}

	console.Runner = runner.NewRunner(ctx, console)
	return console, nil
}

func (console *Console) Name() string {
	return "Console"
}

func (console *Console) Setup() {
	go console.Scan()
}

func (console *Console) Run() {
	<-console.Done()
}

func (console *Console) Cleanup() {
	os.Stdin.Close()
}

func (console *Console) Scan() {
	for console.Scanner.Scan() {
		log.Infof("%s", console.Text())
		fmt.Fprint(os.Stderr, "> ")
	}

	if err := console.Err(); err != nil {
		log.Errorf("Error reading input from console: %s", err)
	}
}
