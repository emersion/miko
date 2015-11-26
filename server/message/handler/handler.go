// Handles incoming messages from remotes.
package handler

import (
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
	"time"
)

// A handler for a specific message
type TypeHandler func(*message.Context, *message.IO) error

// Handles messages from remote
type Handler struct {
	ctx      *message.Context
	handlers map[message.Type]TypeHandler
}

// Handle a message of the specified type
func (h *Handler) Handle(t message.Type, io *message.IO) error {
	if h.ctx.IsServer() {
		if t != message.Types["version"] && io.Version == 0 {
			// Client didn't send his version number
			if err := builder.SendExit(io, message.ExitCodes["client_outdated"]); err != nil {
				return err
			}
			return io.Close()
		}
	}

	if val, ok := h.handlers[t]; ok {
		err := val(h.ctx, io)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown message type: %d", t)
	}

	return nil
}

// Listen to a remote stream
func (h *Handler) Listen(clientIO *message.IO) {
	defer (func() {
		// Will be executed when the connection is closed

		clientIO.Close()

		var session *message.Session
		if h.ctx.IsServer() {
			session = h.ctx.Auth.GetSession(clientIO.Id)
			if session != nil {
				h.ctx.Auth.Logout(clientIO.Id)
				builder.SendPlayerLeft(clientIO.Broadcaster(), h.ctx.Clock.GetRelativeTick(), session.Entity.Id)
			}
		}
		if h.ctx.IsClient() {
			session = h.ctx.Me
		}

		if session != nil && session.Entity != nil {
			// TODO: move this elsewhere
			h.ctx.Entity.Delete(session.Entity.Id, h.ctx.Clock.GetAbsoluteTick())
		}
	})()

	var msgType message.Type
	for {
		err := read(clientIO, &msgType)
		if err == io.EOF {
			log.Println("Connection closed.")
			return
		} else if err != nil {
			//clientIO.Close()
			log.Println("binary.Read failed:", err)
			return
		}

		done := make(chan bool)
		go (func() {
			select {
			case <-done:
				return
			case <-time.After(time.Second * 8):
				log.Println("TIMED OUT:", message.GetTypeName(msgType))
				panic("Message handling timed out")
			}
		})()

		err = h.Handle(msgType, clientIO)
		if err != nil {
			log.Println("Handle failed:", err)
			log.Println("Message type:", msgType)
		}

		done <- true
	}
}

func mergeHandlers(handlersList ...*map[message.Type]TypeHandler) map[message.Type]TypeHandler {
	result := map[message.Type]TypeHandler{}

	for _, handlers := range handlersList {
		for t, handler := range *handlers {
			result[t] = handler
		}
	}

	return result
}

func New(ctx *message.Context) *Handler {
	var handlers map[message.Type]TypeHandler

	if ctx.IsServer() {
		handlers = mergeHandlers(commonHandlers, serverHandlers)
	} else if ctx.IsClient() {
		handlers = mergeHandlers(commonHandlers, clientHandlers)
	} else {
		handlers = *commonHandlers
	}

	return &Handler{
		ctx:      ctx,
		handlers: handlers,
	}
}
