package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/jbhannah/gophermine/pkg/console"
	"github.com/jbhannah/gophermine/pkg/listener"
	"github.com/jbhannah/gophermine/pkg/mc"

	log "github.com/sirupsen/logrus"
)

// RCONServer listens for and handles incoming RCON connections.
type RCONServer struct {
	*listener.Listener
}

// NewRCONServer returns a new RCONServer.
func NewRCONServer(ctx context.Context, addr string) (*RCONServer, error) {
	rcon := &RCONServer{}

	listener, err := listener.NewListener(ctx, rcon, addr)
	if err != nil {
		return nil, err
	}

	rcon.Listener = listener
	return rcon, nil
}

// Name returns the name of the RCON server.
func (rcon *RCONServer) Name() string {
	return "RCON"
}

// HandleConn handles incoming RCON connections.
func (rcon *RCONServer) HandleConn(conn net.Conn) {
	console, err := rcon.newConsole(conn)
	if err != nil {
		log.Errorf("Could not initialize RCON console for %s: %s", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	if _, err := conn.Write([]byte("Password: ")); err != nil {
		log.Warnf("Error prompting for RCON password from %s: %s", conn.RemoteAddr(), err)
	}

	if !console.Scan() || console.Text() != mc.Properties().RCON.Password {
		if _, err := conn.Write([]byte("Incorrect password!\n")); err != nil {
			log.Warnf("Error warning %s about incorrect RCON password: %s", conn.RemoteAddr(), err)
		}

		log.Warnf("Failed RCON authentication attempt from %s", conn.RemoteAddr())

		conn.Close()
		return
	}

	<-console.Start()
	<-console.Stopped()
}

type rconConn struct {
	net.Conn
	*console.Console
}

func (rcon *RCONServer) newConsole(conn net.Conn) (*rconConn, error) {
	rc := &rconConn{
		Conn: conn,
	}

	name := fmt.Sprintf("%s RCON console", conn.RemoteAddr().String())
	writer := io.MultiWriter(rc, log.StandardLogger().WriterLevel(log.DebugLevel))

	con, err := console.NewConsole(rcon.Context, name, rc, writer)
	if err != nil {
		return nil, err
	}

	rc.Console = con
	return rc, nil
}

// Read reads incoming data from the connection and stops the console if the
// connection is closed.
func (conn *rconConn) Read(p []byte) (int, error) {
	n, err := conn.Conn.Read(p)

	if err == io.EOF {
		defer conn.Console.Stop()
	}

	return n, err
}
