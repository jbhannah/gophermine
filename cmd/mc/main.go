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

// DefaultAddr is the default server address
const DefaultAddr = ":25565"

func main() {
	log.Debug("Starting gophermine")

	ctx := context.Background()
	server := server.NewServer(ctx, DefaultAddr)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(server, sigs)

	server.Start()
	log.Debug("Started gophermine")

	<-server.Stopped()

	close(sigs)
	log.Debug("Stopped gophermine")
}

func handleSigs(server *server.Server, sigs <-chan os.Signal) {
	sig := <-sigs

	print("\r")
	log.Debug(fmt.Sprintf("Received %s signal", sig))
	log.Debug("Stopping gophermine")

	server.Stop()
}
