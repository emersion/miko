// Provides a TCP server, with TLS.
package server

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"git.emersion.fr/saucisse-royale/miko.git/server/crypto"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
)

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

	address string           // Address to open connection, e.g. localhost:9999
	Joins   chan *message.IO // Channel for new connections
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
	c.Server.ios[c.id] = nil

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

	r := bufio.NewReader(c.conn)
	w := bufio.NewWriter(c)
	io := message.NewIO(c.id, r, w, c, s)

	s.ios = append(s.ios, io)
	c.Server.Joins <- io
}

// Start network server
func (s *Server) Listen() {
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

	go (func() {
		for {
			ready := 0
			for _, io := range s.ios {
				if io == nil || io.State != message.Ready {
					continue
				}

				ready++
			}
			log.Println("Users:", len(s.ios), "ready:", ready)
			time.Sleep(time.Second)
		}
	})()

	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		s.newClient(conn)
	}
}

func (s *Server) Lock() {
	var wg sync.WaitGroup

	for _, io := range s.ios {
		if io == nil || io.State != message.Ready {
			continue
		}

		wg.Add(1)

		go (func(io *message.IO) {
			defer wg.Done()

			io.Lock()
		})(io)
	}

	wg.Wait()
}

func (s *Server) Unlock() {
	for _, io := range s.ios {
		if io == nil || io.State != message.Ready {
			continue
		}

		io.Unlock()
	}
}

// Broadcast a message to all clients
func (s *Server) Write(data []byte) (n int, err error) {
	N := 0

	for _, io := range s.ios {
		if io == nil || io.State != message.Ready {
			continue
		}

		n, err := io.Write(data)
		if err != nil {
			log.Println("Warning: error while broadcasting data:", err)
		}

		N += n
	}

	return N, nil
}

// Creates new tcp server instance
func New(address string) *Server {
	log.Println("Creating server with address", address)
	server := &Server{
		address: address,
		Joins:   make(chan *message.IO),
	}

	return server
}
