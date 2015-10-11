package server

import (
	"log"
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

	clientIO := &message.IO{
		Reader: c.conn,
		Writer: c.conn,
		BroadcastWriter: c.Server,
		Id: c.id,
	}

	c.Server.handler.Listen(clientIO)
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
		n, err = c.conn.Write(msg)
		if err != nil {
			return N, err
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
