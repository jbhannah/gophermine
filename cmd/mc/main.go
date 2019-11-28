package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbhannah/gophermine/pkg/console"
	"github.com/jbhannah/gophermine/pkg/mc"

	"github.com/jbhannah/gophermine/internal/pkg/server"
	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
)

const (
	// Version is the Gophermine release version
	Version = "0.0.1"

	// MCVersion is the Minecraft™: Java Edition version that this release of
	// Gophermine is compatible with.
	MCVersion = "1.14.4"

	// MCProtocolVersion is the Minecraft protocol version number that this
	// release of Gophermine is compatible with.
	MCProtocolVersion = 498
)

var (
	help    bool
	verbose bool
	version bool
)

func init() {
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000Z-07:00",
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		log.SetFormatter(&console.TermFormatter{
			TextFormatter: formatter,
		})
	} else {
		log.SetFormatter(formatter)
	}
}

func main() {
	parseFlags()

	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	log.Info("Starting Gophermine")
	startTime := time.Now()

	if err := mc.CheckEULA(); err != nil {
		return err
	}

	if err := mc.LoadProperties(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(cancel, sigs)

	server, err := server.NewServer(ctx)
	if err != nil {
		return err
	}

	<-server.Start()
	log.Infof("Started Gophermine in %s", time.Since(startTime))

	<-server.Stopped()
	log.Infof("Stopped Gophermine after %s", time.Since(startTime))

	return nil
}

func handleSigs(cancel context.CancelFunc, sigs <-chan os.Signal) {
	defer cancel()
	sig := <-sigs

	print("\r")
	log.Warn(fmt.Sprintf("Received %s signal", sig))
	log.Warn("Stopping Gophermine")
}
