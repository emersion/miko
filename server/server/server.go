package server

import (
	"bufio"
	"log"
	"net"
	"encoding/binary"
	"fmt"
)

// Client holds info about connection
type Client struct {
	conn net.Conn
	Server *server
	incoming chan string // Channel for incoming data from client
}

// TCP server
type server struct {
	clients []*Client
	address string // Address to open connection: localhost:9999
	joins chan net.Conn // Channel for new connections
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)

	var msg_type uint16
	for {
		err := binary.Read(reader, binary.BigEndian, &msg_type)
		if err != nil {
			c.conn.Close()
			log.Println("binary.Read failed:", err)
			return
		}
		fmt.Print(msg_type)
	}
}

// Send message to client
func (c *Client) Send(msg []byte) error {
	_, err := c.conn.Write(msg)
	return err
}

// Creates new Client instance and starts listening
func (s *server) newClient(conn net.Conn) {
	client := &Client{
		conn:   conn,
		Server: s,
	}
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

// Start network server
func (s *server) Listen() {
	go s.listenChannels()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		s.joins <- conn
	}
}

// Creates new tcp server instance
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
		joins:   make(chan net.Conn),
	}

	return server
}
