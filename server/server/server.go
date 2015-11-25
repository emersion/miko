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
	conn   net.Conn
	Server *Server
	id     int
}

// TCP server
type Server struct {
	clients []*Client
	ios     []*message.IO

	brd    []chan []byte
	brdEnd []chan bool
	locked bool

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

	s.brd = append(s.brd, make(chan []byte))
	s.brdEnd = append(s.brdEnd, make(chan bool))
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
	if s.locked {
		return
	}
	s.locked = true

	for i, io := range s.ios {
		if io == nil {
			continue
		}

		buffer := [][]byte{}
		locked := false
		finished := false

		go (func() {
			io.Lock()

			for _, data := range buffer {
				io.Write(data)
			}

			locked = true
			if finished {
				io.Unlock()
			}
		})()

		go (func() {
			for {
				select {
				case data := <-s.brd[i]:
					if !locked {
						buffer = append(buffer, data)
					} else {
						io.Write(data)
					}
				case <-s.brdEnd[i]:
					finished = true

					if locked {
						io.Unlock()
					}
					return
				}
			}
		})()
	}
}

func (s *Server) Unlock() {
	if !s.locked {
		return
	}
	s.locked = false

	for i, io := range s.ios {
		if io == nil {
			continue
		}

		select {
		case s.brdEnd[i] <- true:
		default:
		}
	}
}

// Broadcast a message to all clients
func (s *Server) Write(data []byte) (n int, err error) {
	for i, io := range s.ios {
		if io == nil {
			continue
		}

		if s.locked {
			select {
			case s.brd[i] <- data:
			default:
			}
		} else {
			io.Write(data)
		}
	}
	return len(data), nil
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
