// Provides a TCP server, with TLS.
package server

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"

	"git.emersion.fr/saucisse-royale/miko.git/server/crypto"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

type broadcaster struct {
	id     int
	server *Server
}

func (brd *broadcaster) Write(data []byte) (n int, err error) {
	return brd.server.broadcast(data, brd.id)
}

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *Server
	id     int
}

// TCP server
type Server struct {
	clients []*Client
	ios     []*message.IO
	address string        // Address to open connection, e.g. localhost:9999
	joins   chan net.Conn // Channel for new connections
	Joins   chan *message.IO
}

func (c *Client) Write(data []byte) (n int, err error) {
	return c.conn.Write(data)
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
func (s *Server) newClient(conn net.Conn) {
	c := &Client{
		conn:   conn,
		Server: s,
		id:     len(s.clients),
	}
	s.clients = append(s.clients, c)

	log.Println("New client:", c.id)

	reader := bufio.NewReader(c.conn)
	brd := &broadcaster{c.id, c.Server}
	io := message.NewIO(c.id, reader, c, brd)

	s.ios = append(s.ios, io)
	c.Server.Joins <- io
}

// Listens new connections channel and creating new client
func (s *Server) listenChannels() {
	for {
		select {
		case conn := <-s.joins:
			s.newClient(conn)
		}
	}
}

// Start network server
func (s *Server) Listen() {
	go s.listenChannels()

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

func (s *Server) broadcast(data []byte, from int) (n int, err error) {
	N := 0
	for _, c := range s.clients {
		if c == nil {
			continue
		}
		if from != c.id {

		}

		n, err = c.conn.Write(data)
		if err != nil {
			log.Println("Error broadcasting message:", err)
		}
		N += n
	}
	return N, nil
}

// Broadcast a message to all clients
func (s *Server) Write(data []byte) (n int, err error) {
	return s.broadcast(data, -1)
}

// Creates new tcp server instance
func New(address string) *Server {
	log.Println("Creating server with address", address)
	server := &Server{
		address: address,
		joins:   make(chan net.Conn),
		Joins:   make(chan *message.IO),
	}

	return server
}
