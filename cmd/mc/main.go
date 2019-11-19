package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbhannah/gophermine/internal/pkg/server"
)

func main() {
	println("Starting gophermine")

	ctx, cancel := context.WithCancel(context.Background())
	server := server.NewServer(ctx, cancel)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(server, sigs)

	server.Start()
	println("Started gophermine")

	<-ctx.Done()
	println("Stopped gophermine")
}

func handleSigs(server *server.Server, sigs <-chan os.Signal) {
	sig := <-sigs

	print("\r")
	println(fmt.Sprintf("Received %s signal", sig))
	println("Stopping gophermine")

	server.Stop()
}
