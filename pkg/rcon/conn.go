package rcon

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/jbhannah/gophermine/pkg/console"
	"github.com/jbhannah/gophermine/pkg/mc"

	log "github.com/sirupsen/logrus"
)

// Conn is an open RCON connection for receiving and responding to commands.
type Conn struct {
	net.Conn
	*console.Console
	response chan []byte
	errors   chan error
}

// NewConn authenticates and opens a console for a new RCON connection.
func NewConn(ctx context.Context, conn net.Conn) (*Conn, error) {
	c := &Conn{
		Conn:     conn,
		response: make(chan []byte),
		errors:   make(chan error),
	}

	name := fmt.Sprintf("%s RCON console", conn.RemoteAddr())
	if con, err := console.NewConsole(ctx, name, c, c); err != nil {
		return nil, err
	} else {
		c.Console = con
	}

	return c, c.AcceptLogin()
}

// ReadPacket builds a Packet from an incoming RCON packet.
func (conn *Conn) ReadPacket() (*Packet, error) {
	packet := &Packet{}
	reader := conn.Conn

	if err := binary.Read(reader, binary.LittleEndian, &packet.Length); err != nil {
		return nil, err
	}

	if packet.Length < 10 {
		return nil, fmt.Errorf("RCON packet length of %d is too short", packet.Length)
	}

	buf := make([]byte, packet.Length)
	if err := binary.Read(reader, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}

	packet.RequestID = int32(binary.LittleEndian.Uint32(buf[:4]))
	packet.Type = PacketType(binary.LittleEndian.Uint32(buf[4:8]))
	packet.Payload = buf[8 : packet.Length-2]

	return packet, nil
}

// Read handles incoming RCON messages. If an EOF is encountered, the connection
// stops its running console. If no error is encountered, a response with the
// same request ID as the incoming packet is built and sent once a response is
// written to the connection's console.
func (conn *Conn) Read(p []byte) (int, error) {
	packet, err := conn.ReadPacket()

	if err == io.EOF {
		defer conn.Console.Stop()
		return 0, err
	} else if err != nil {
		return 0, err
	}

	if err := packet.ValidateType(Command); err != nil {
		return 0, err
	}

	go func() {
		_, err := conn.WritePacket(packet.RequestID, Response, <-conn.response)
		conn.errors <- err
	}()

	buf := new(bytes.Buffer)
	buf.Write(packet.Payload)
	buf.WriteString("\n")

	n := copy(p, buf.Bytes())
	return n, nil
}

// WritePacket builds an RCON packet and writes it to the underlying network
// connection.
func (conn *Conn) WritePacket(id int32, pt PacketType, payload []byte) (int, error) {
	packet := NewPacket(id, pt, payload)

	bytes, err := packet.Bytes()
	if err != nil {
		return 0, err
	}

	return conn.Conn.Write(bytes)
}

// Write logs responses to commands sent via the RCON connection, passes the
// response payload to response channel, and returns any errors encountered by
// the writing of the packet to the underlying network connection.
func (conn *Conn) Write(p []byte) (int, error) {
	log.Info(string(p))
	conn.response <- p
	return len(p), <-conn.errors
}

// AcceptLogin reads an AuthRequest RCON packet and validates its payload
// against the RCON password configured in the server.properties file.
func (conn *Conn) AcceptLogin() error {
	packet, err := conn.ReadPacket()
	if err != nil {
		return err
	}

	if err := packet.ValidateType(AuthRequest); err != nil {
		return err
	}

	if string(packet.Payload) != mc.Properties().RCON.Password {
		if _, err := conn.WritePacket(-1, AuthResponse, make([]byte, 0)); err != nil {
			log.Errorf("Error responding to invalid RCON authentication request: %s", err)
		}

		return fmt.Errorf("Incorrect password")
	}

	_, err = conn.WritePacket(packet.RequestID, AuthResponse, make([]byte, 0))
	return err
}
