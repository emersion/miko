package timeserver

import (
	"bytes"
	"encoding/binary"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"net"
	"time"
)

type timestampReader struct {
	conn   *net.UDPConn
	buffer []byte
	reader *bytes.Reader
}

func (r *timestampReader) Read() (timestamp uint64, addr *net.UDPAddr, err error) {
	r.reader.Reset(r.buffer)
	_, addr, err = r.conn.ReadFromUDP(r.buffer)
	if err != nil {
		return
	}

	err = binary.Read(r.reader, binary.BigEndian, &timestamp)
	return
}

func newTimestampReader(conn *net.UDPConn) *timestampReader {
	buf := make([]byte, 8)
	return &timestampReader{
		conn:   conn,
		buffer: buf,
		reader: bytes.NewReader(buf),
	}
}

type Server struct {
	addr    *net.UDPAddr
	conn    *net.UDPConn
	clients map[string]*Client

	Joins chan *Client
}

type Client struct {
	addr     *net.UDPAddr
	accepted bool
	received uint64
}

func (s *Server) replyTo(client *Client) error {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, client.received)
	if err != nil {
		return err
	}

	now := message.TimeToTimestamp(time.Now())
	err = binary.Write(b, binary.BigEndian, now)
	if err != nil {
		return err
	}

	_, err = s.conn.WriteToUDP(b.Bytes(), client.addr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Accept(c *Client) error {
	if c.accepted {
		return nil // Already accepted
	}

	c.accepted = true

	return s.replyTo(c)
}

func (s *Server) Reject(c *Client) {
	if !c.accepted {
		return
	}

	// Do not do anything
}

func (s *Server) Listen() error {
	var err error
	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}

	defer s.conn.Close()

	r := newTimestampReader(s.conn)
	for {
		timestamp, addr, err := r.Read()
		if err != nil {
			return err
		}

		if client, ok := s.clients[addr.String()]; ok {
			client.received = timestamp
			if client.accepted {
				s.replyTo(client)
			}
		} else {
			client := &Client{
				addr:     addr,
				received: timestamp,
			}
			s.clients[addr.String()] = client
			s.Joins <- client
		}
	}
}

func (s *Server) Port() int {
	return s.addr.Port
}

func New(addrStr string) *Server {
	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		panic("Cannot resolve UDP address")
	}

	return &Server{
		addr:  addr,
		Joins: make(chan *Client),
		clients: make(map[string]*Client),
	}
}
