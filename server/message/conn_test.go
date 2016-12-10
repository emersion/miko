package message_test

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
)

type mockServer struct {
	conns []*message.Conn
}

func (s *mockServer) Write(write message.WriteFunc) error {
	for _, c := range s.conns {
		c.Write(write)
	}
	return nil
}

func TestConn_one(t *testing.T) {
	conn := newMockConn()

	// One message sent over one conn
	conn.Write(builder.SendPing)
	msgType := handler.ReadType(conn)
	if msgType != message.Types["ping"] {
		t.Error("Sent ping, but didn't received it")
	}

	// Multiple messages sent over one io
	chatCount := 10
	chatMsg := "Hello World! What's up bro?"
	var wg sync.WaitGroup
	for i := 0; i < chatCount; i++ {
		wg.Add(1)
		go (func() {
			defer wg.Done()
			conn.Write(func (w io.Writer) error {
				return builder.SendChatSend(w, chatMsg)
			})
		})()
	}
	wg.Wait()

	for i := 0; i < chatCount; i++ {
		msgType := handler.ReadType(conn)
		if msgType != message.Types["chat_send"] {
			t.Error("Sent chat_send, but didn't received it at iteration", i)
		}
		receivedMsg := handler.ReadChatSend(conn)
		if receivedMsg != chatMsg {
			t.Error("Bad received chat message content at iteration", i)
		}
	}
}

func TestConn_multiple(t *testing.T) {
	s := &mockServer{}

	var buffers []*bytes.Buffer

	for i := 0; i < 10; i++ {
		// TODO: use newMockConn()
		b := &bytes.Buffer{}
		conn := message.NewConn(0, b, b, s.Write)
		s.conns = append(s.conns, conn)
		buffers = append(buffers, b)
	}

	// One message sent over all ios
	s.Write(func (w io.Writer) error {
		return builder.SendPing(w)
	})

	for i, conn := range s.conns {
		msgType := handler.ReadType(conn)
		if msgType != message.Types["ping"] {
			t.Error("Sent ping, but didn't received it over io", i)
		}
	}

	// Multiple messages sent over multiple conns
	chatCount := 10
	chatMsg := "Hello World! What's up bro?"
	var wg sync.WaitGroup
	for i := 0; i < chatCount; i++ {
		wg.Add(1)
		go (func() {
			defer wg.Done()

			s.Write(func (w io.Writer) error {
				return builder.SendChatSend(w, chatMsg)
			})
		})()
	}
	wg.Wait()

	for _, conn := range s.conns {
		for i := 0; i < chatCount; i++ {
			msgType := handler.ReadType(conn)
			if msgType != message.Types["chat_send"] {
				t.Error("Sent chat_send, but didn't received it at iteration", i)
			}
			receivedMsg := handler.ReadChatSend(conn)
			if receivedMsg != chatMsg {
				t.Error("Bad received chat message content at iteration", i)
			}
		}
	}
}
