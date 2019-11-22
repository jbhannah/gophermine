package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbhannah/gophermine/internal/pkg/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000Z-07:00",
	})
	log.SetLevel(log.DebugLevel)
}

// DefaultRCONAddr is the default server address
const DefaultRCONAddr = ":25566"

func main() {
	log.Info("Starting Gophermine")

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(cancel, sigs)

	server, err := server.NewServer(ctx, DefaultRCONAddr)
	if err != nil {
		log.Fatal(err)
	}

	server.Start()
	log.Info("Started Gophermine")

	<-server.Stopped()
	log.Info("Stopped Gophermine")
}

func handleSigs(cancel context.CancelFunc, sigs <-chan os.Signal) {
	defer cancel()
	sig := <-sigs

	print("\r")
	log.Warn(fmt.Sprintf("Received %s signal", sig))
	log.Warn("Stopping Gophermine")
}
