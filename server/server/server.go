// Provides a TCP server, with TLS.
package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
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

	wc := struct{
		io.Writer
		io.Closer
	}{w, conn}

	c := &Client{
		Conn: message.NewConn(len(s.clients), r, wc, s.Write),
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
func (s *Server) Write(write message.WriteFunc) error {
	// Write to a buffer
	var buf bytes.Buffer
	if err := write(&buf); err != nil {
		return err
	}

	n := 0
	done := make(chan error)
	for _, c := range s.clients {
		if c.State != message.Ready {
			continue
		}

		go func() {
			done <- c.Write(func(w io.Writer) error {
				_, err := w.Write(buf.Bytes())
				return err
			})
		}()

		n++
	}

	// TODO: handle per-conn errors
	for i := 0; i < n; i++ {
		<-done
	}
	return nil
}

// Creates new tcp server instance
func New(address string) *Server {
	return &Server{
		address: address,
		Joins:   make(chan *Client),
	}
}
