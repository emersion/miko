package timeserver

import (
	"bytes"
	"encoding/binary"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"net"
	"time"
)

const interval = time.Millisecond * 100
const timeout = time.Second * 15

type Server struct {
	addr    *net.UDPAddr
	conn    *net.UDPConn
	clients map[string]*Client

	Joins chan *Client
}

type Client struct {
	addr     *net.UDPAddr
	accepted bool
	stop     chan bool
	reset    chan bool
}

func (s *Server) Accept(c *Client) error {
	if c.accepted {
		return nil // Already accepted
	}

	c.accepted = true

	b := new(bytes.Buffer)
	ticker := time.NewTicker(interval)
	delay := time.NewTimer(timeout)

	defer (func() {
		c.accepted = false
		ticker.Stop()
		delay.Stop()
	})()

	for {
		select {
		case t := <-ticker.C:
			b.Reset()
			err := binary.Write(b, binary.BigEndian, message.TimeToTimestamp(t))
			if err != nil {
				return err
			}
			_, err = s.conn.WriteToUDP(b.Bytes(), c.addr)
			if err != nil {
				return err
			}
		case <-c.stop:
			return nil
		case <-c.reset:
			delay.Reset(timeout)
		case <-delay.C:
			return nil
		}
	}
}

func (s *Server) Reject(c *Client) {
	if !c.accepted {
		return
	}

	c.stop <- true
}

func (s *Server) Listen() error {
	var err error
	s.conn, err = net.ListenUDP("udp", s.addr)
	if err != nil {
		return err
	}

	defer s.conn.Close()

	buf := make([]byte, 1) // We should receive only empty packets
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// Only process empty packets
		if n != 0 {
			continue
		}

		if client, ok := s.clients[addr.String()]; ok {
			client.reset <- true
		} else {
			client := &Client{
				addr:  addr,
				stop:  make(chan bool, 1),
				reset: make(chan bool, 1),
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
	}
}
