package server

import (
	"context"
	"net"

	"github.com/jbhannah/gophermine/pkg/listener"
	"github.com/jbhannah/gophermine/pkg/rcon"

	log "github.com/sirupsen/logrus"
)

// RCONServer listens for and handles incoming RCON connections.
type RCONServer struct {
	*listener.Listener
}

// NewRCONServer returns a new RCONServer.
func NewRCONServer(ctx context.Context, addr string) (*RCONServer, error) {
	srv := &RCONServer{}

	listener, err := listener.NewListener(ctx, srv, addr)
	if err != nil {
		return nil, err
	}

	srv.Listener = listener
	return srv, nil
}

// Name returns the name of the RCON server.
func (srv *RCONServer) Name() string {
	return "RCON"
}

// HandleConn handles incoming RCON connections.
func (srv *RCONServer) HandleConn(conn net.Conn) {
	rconn, err := rcon.NewConn(srv.Context, conn)
	if err != nil {
		log.Errorf("Could not initialize RCON console for %s: %s", conn.RemoteAddr(), err)
		conn.Close()
		return
	}

	<-rconn.Start()
	<-rconn.Stopped()
}
