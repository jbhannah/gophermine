package server

import (
	"context"
	"io"
	"net"

	"github.com/jbhannah/gophermine/pkg/listener"
	log "github.com/sirupsen/logrus"
)

// MCServer listens for and handles incoming Minecraft client connections.
type MCServer struct {
	*listener.Listener
}

// NewMCServer returns a new MCServer.
func NewMCServer(ctx context.Context, addr string) (*MCServer, error) {
	mc := &MCServer{}

	listener, err := listener.NewListener(ctx, mc, addr)
	if err != nil {
		return nil, err
	}

	mc.Listener = listener
	return mc, nil
}

// Name returns the name of the Minecraft server.
func (mc *MCServer) Name() string {
	return "Minecraft"
}

// HandleConn handles incoming Minecraft connections.
func (mc *MCServer) HandleConn(conn net.Conn) {
	if _, err := io.Copy(conn, conn); err != nil {
		log.Errorf("Error in connection: %s", err)
	}
}
