package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbhannah/gophermine/pkg/mc"

	"github.com/jbhannah/gophermine/internal/pkg/server"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	// Version is the Gophermine release version
	Version = "0.0.1"

	// MCVersion is the Minecraft™: Java Edition version that this release of
	// Gophermine is compatible with.
	MCVersion = "1.14.4"

	// MCProtocolVersion is the Minecraft protocol version number that this release
	// of Gophermine is compatible with.
	MCProtocolVersion = 498
)

var (
	help       bool
	rconPort   int
	serverPort int
	verbose    bool
	version    bool
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000000Z-07:00",
	})

	flag.BoolVarP(&help, "help", "h", false, "show this help message")
	flag.IntVar(&rconPort, "rcon-port", mc.RCONPort, "port to listen on for RCON connections")
	flag.IntVarP(&serverPort, "port", "p", mc.ServerPort, "port to listen on for Minecraft client connections")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVarP(&version, "version", "v", false, "print the version")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	if help {
		printVersion()
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(0)
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Enabled verbose logging")
	}

	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	log.Info("Starting Gophermine")

	config, err := mc.NewServerConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSigs(cancel, sigs)

	server, err := server.NewServer(ctx, config.ServerAddr(), config.RCONAddr())
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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage of mc:")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "Gophermine %s (Java Edition version %s, protocol %d)\n", Version, MCVersion, MCProtocolVersion)
}
