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

// Client holds info about connection
type Client struct {
	*message.Conn
	Server *Server
}

func (c *Client) Close() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}

	c.Server.clients = append(c.Server.clients[:c.Id], c.Server.clients[c.Id+1:]...)
	return nil
}

type Server struct {
	clients []*Client

	address string       // Address to open connection, e.g. localhost:9999
	Joins   chan *Client // Channel for new connections
}

// Creates new Client instance and starts listening
func (s *Server) newClient(conn net.Conn) {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	c := &Client{
		Conn: message.NewConn(len(s.clients), r, w),
		Server: s,
	}

	s.clients = append(s.clients, c)
	c.Server.Joins <- c
}

// Start network server
func (s *Server) Listen() error {
	tlsConfig, err := crypto.GetServerTlsConfig()
	if err != nil {
		return err
	}

	var listener net.Listener
	if tlsConfig != nil {
		listener, err = tls.Listen("tcp", s.address, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", s.address)
		log.Println("Warning: creating a non-TLS insecure server")
	}
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		s.newClient(conn)
	}

	return nil
}

// Broadcast a message to all clients
func (s *Server) Write(data []byte) (int, error) {
	N := 0

	for _, c := range s.clients {
		if c.State != message.Ready {
			continue
		}

		n, err := c.Write(data)
		if err != nil {
			log.Println("Warning: error while broadcasting data:", err)
		}

		N += n
	}

	return N, nil
}

// Creates new tcp server instance
func New(address string) *Server {
	return &Server{
		address: address,
		Joins:   make(chan *Client),
	}
}
