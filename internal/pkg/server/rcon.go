package server

import (
	"context"
	"io"
	"net"

	"github.com/jbhannah/gophermine/pkg/listener"
	log "github.com/sirupsen/logrus"
)

// RCONServer listens for and handles incoming RCON connections.
type RCONServer struct {
	*listener.Listener
}

// NewRCONServer returns a new RCONServer.
func NewRCONServer(ctx context.Context, addr string) *RCONServer {
	rcon := &RCONServer{}
	rcon.Listener = listener.NewListener(ctx, rcon, addr)

	return rcon
}

// HandleConn handles incoming RCON connections.
func (rcon *RCONServer) HandleConn(conn net.Conn) {
	log.Debugf("Accepted connection from %s", conn.RemoteAddr())

	if _, err := io.Copy(conn, conn); err != nil {
		log.Errorf("Error in connection: %s", err)
	}

	log.Debugf("Closed connection from %s", conn.RemoteAddr())
}
