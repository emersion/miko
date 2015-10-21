package server

import (
	"log"
	"bufio"
	"net/http"
	"golang.org/x/net/websocket"

	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/handler"
)

// Client holds info about connection
type Client struct {
	conn *websocket.Conn
	Server *server
	incoming chan string // Channel for incoming data from client
	id int
}

// TCP server
type server struct {
	clients []*Client
	address string // Address to open connection, e.g. localhost:9999
	handler *handler.Handler
}

// Read client data from channel
func (c *Client) listen() {
	log.Println("New client:", c.id)

	reader := bufio.NewReader(c.conn)

	// Send binary frames
	// See http://grokbase.com/t/gg/golang-nuts/1314ee50mh/go-nuts-sending-binary-websocket-frames#20130104bmvbvtymkitzju2nhnaytrrzs4
	c.conn.PayloadType = 0x2

	clientIO := &message.IO{
		Reader: reader,
		Writer: c.conn,
		BroadcastWriter: c.Server,
		Id: c.id,
	}

	defer c.Close()

	c.Server.handler.Listen(clientIO)
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.Server.clients[c.id] = nil

	return nil
}

func (s *server) newClient(conn *websocket.Conn) {
	client := &Client{
		conn: conn,
		Server: s,
		id: len(s.clients),
	}
	s.clients = append(s.clients, client)
	client.listen()
}

// Start network server
func (s *server) Listen() {
	http.Handle("/socket", websocket.Handler(s.newClient))
	http.Handle("/", http.FileServer(http.Dir("public")))

	err := http.ListenAndServe(s.address, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
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

// Creates new http server instance
func New(address string, ctx *message.Context) *server {
	log.Println("Creating HTTP server with address", address)
	server := &server{
		address: address,
		handler: handler.New(ctx),
	}

	return server
}
