package message_test

import (
	"bytes"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/handler"
	"sync"
	"testing"
)

type mockServer struct {
	ios []*message.IO
}

func (s *mockServer) Write(p []byte) (n int, err error) {
	for _, io := range s.ios {
		io.Write(p)
	}
	return
}

func (s *mockServer) Lock() {
	for _, io := range s.ios {
		io.Lock()
	}
}

func (s *mockServer) Unlock() {
	for _, io := range s.ios {
		io.Unlock()
	}
}

func TestIO_one(t *testing.T) {
	io := newMockIO()

	// One message sent over one io
	builder.SendPing(io)
	msgType := handler.ReadType(io)
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
			builder.SendChatSend(io, chatMsg)
		})()
	}
	wg.Wait()

	for i := 0; i < chatCount; i++ {
		msgType := handler.ReadType(io)
		if msgType != message.Types["chat_send"] {
			t.Error("Sent chat_send, but didn't received it at iteration", i)
		}
		receivedMsg := handler.ReadChatSend(io)
		if receivedMsg != chatMsg {
			t.Error("Bad received chat message content at iteration", i)
		}
	}
}

func TestIO_multiple(t *testing.T) {
	s := &mockServer{}

	var buffers []*bytes.Buffer

	for i := 0; i < 10; i++ {
		// TODO: use newMockIO()
		b := &bytes.Buffer{}
		c := &mockCloser{}
		io := message.NewIO(0, b, b, c, s)
		s.ios = append(s.ios, io)
		buffers = append(buffers, b)
	}

	// One message sent over all ios
	builder.SendPing(s)

	for i, io := range s.ios {
		msgType := handler.ReadType(io)
		if msgType != message.Types["ping"] {
			t.Error("Sent ping, but didn't received it over io", i)
		}
	}

	// Multiple messages sent over multiple ios
	chatCount := 10
	chatMsg := "Hello World! What's up bro?"
	var wg sync.WaitGroup
	for i := 0; i < chatCount; i++ {
		wg.Add(1)
		go (func() {
			defer wg.Done()
			builder.SendChatSend(s, chatMsg)
		})()
	}
	wg.Wait()

	for _, io := range s.ios {
		for i := 0; i < chatCount; i++ {
			msgType := handler.ReadType(io)
			if msgType != message.Types["chat_send"] {
				t.Error("Sent chat_send, but didn't received it at iteration", i)
			}
			receivedMsg := handler.ReadChatSend(io)
			if receivedMsg != chatMsg {
				t.Error("Bad received chat message content at iteration", i)
			}
		}
	}
}
