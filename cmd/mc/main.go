package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbhannah/gophermine/internal/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000Z-07:00",
	})
	log.SetLevel(log.DebugLevel)
}

// DefaultMCAddr is the default Minecraft server address
const DefaultMCAddr = ":25565"

// DefaultRCONAddr is the default RCON server address
const DefaultRCONAddr = ":25566"

func main() {
	app := &cli.App{
		Name:    "Gophermine",
		Usage:   "A (someday) Minecraft™: Java Edition compatible server written in Go.",
		Version: "0.0.1",
		Action:  start,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	log.Info("Starting Gophermine")

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(cancel, sigs)

	server, err := server.NewServer(ctx, DefaultMCAddr, DefaultRCONAddr)
	if err != nil {
		return err
	}

	server.Start()
	log.Info("Started Gophermine")

	<-server.Stopped()
	log.Info("Stopped Gophermine")

	return nil
}

func handleSigs(cancel context.CancelFunc, sigs <-chan os.Signal) {
	defer cancel()
	sig := <-sigs

	print("\r")
	log.Warn(fmt.Sprintf("Received %s signal", sig))
	log.Warn("Stopping Gophermine")
}
