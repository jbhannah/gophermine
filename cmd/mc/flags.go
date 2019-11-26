package main

import (
	"fmt"
	"os"

	"github.com/jbhannah/gophermine/pkg/mc"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func init() {
	flag.BoolVarP(&help, "help", "h", false, "show this help message")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.BoolVarP(&version, "version", "v", false, "print the version")

	flag.IntP("port", "p", mc.ServerPort, "port to listen on for Minecraft client connections")
	if err := mc.Properties().BindPFlag("server-port", flag.Lookup("port")); err != nil {
		log.Fatal(err)
	}

	flag.Usage = usage
}

func parseFlags() {
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
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage of mc:")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "Gophermine %s (Java Edition version %s, protocol %d)\n", Version, MCVersion, MCProtocolVersion)
}
