// Provides a TCP server, with TLS.
package server

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/crypto"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

const broadcastInterval time.Duration = time.Millisecond * 200

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
	id     int
}

// TCP server
type server struct {
	clients []*Client
	ctx     *message.Context
	address string        // Address to open connection, e.g. localhost:9999
	joins   chan net.Conn // Channel for new connections
	handler *handler.Handler
}

// Read client data from channel
func (c *Client) listen() {
	log.Println("New client:", c.id)

	defer c.Close()

	reader := bufio.NewReader(c.conn)
	io := message.NewIO(c.id, reader, c.conn, c.Server)
	c.Server.handler.Listen(io)
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.Server.clients[c.id] = nil

	return nil
}

// Creates new Client instance and starts listening
func (s *server) newClient(conn net.Conn) {
	client := &Client{
		conn:   conn,
		Server: s,
		id:     len(s.clients),
	}
	s.clients = append(s.clients, client)
	go client.listen()
}

// Listens new connections channel and creating new client
func (s *server) listenChannels() {
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		}
	}
}

func (s *server) broadcastChanges() {
	for {
		log.Println("Broadcast")
		if s.ctx.Entity.IsDirty() {
			err := builder.SendEntitiesDiffToClients(s, s.ctx.Clock.GetRelativeTick(), s.ctx.Entity.Flush())
			if err != nil {
				log.Println("Cannot broadcast entities diff:", err)
			}
		}
		if s.ctx.Action.IsDirty() {
			err := builder.SendActionsDone(s, s.ctx.Clock.GetRelativeTick(), s.ctx.Action.Flush())
			if err != nil {
				log.Println("Cannot broadcast actions:", err)
			}
		}

		time.Sleep(broadcastInterval)
	}
}

// Start network server
func (s *server) Listen() {
	go s.listenChannels()
	go s.broadcastChanges()

	tlsConfig, err := crypto.GetServerTlsConfig()
	if err != nil {
		log.Println("Warning: could not get TLS config")
	}

	var listener net.Listener
	if tlsConfig != nil {
		listener, err = tls.Listen("tcp", s.address, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", s.address)
		log.Println("Warning: creating a non-TLS insecure server")
	}

	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		s.joins <- conn
	}
}

// Broadcast a message to all clients
func (s *server) Write(msg []byte) (n int, err error) {
	N := 0
	for _, c := range s.clients {
		if c == nil {
			continue
		}

		n, err = c.conn.Write(msg)
		if err != nil {
			log.Println("Error broadcasting message:", err)
		}
		N += n
	}
	return N, nil
}

// Creates new tcp server instance
func New(address string, ctx *message.Context) *server {
	log.Println("Creating server with address", address)
	server := &server{
		ctx:     ctx,
		address: address,
		joins:   make(chan net.Conn),
		handler: handler.New(ctx),
	}

	return server
}
