package handler

import (
	"fmt"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
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
	log.Println("Received:", message.GetTypeName(t))

	if h.ctx.IsServer() {
		if t != message.Types["version"] && io.Version == 0 {
			// Client didn't send his version number
			if err := builder.SendExit(io.Writer, message.ExitCodes["client_outdated"]); err != nil {
				return err
			}
			return io.Writer.Close()
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

	if h.ctx.IsServer() {
		// No errors, send updates
		if h.ctx.Entity.IsDirty() {
			err := builder.SendEntitiesDiffToClients(io.BroadcastWriter, h.ctx.Clock.GetTick(), h.ctx.Entity.Flush())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Listen to a remote stream
func (h *Handler) Listen(clientIO *message.IO) {
	defer (func() {
		// Will be executed when the connection is closed

		var session *message.Session
		if h.ctx.IsServer() {
			session = h.ctx.Auth.GetSession(clientIO.Id)
			if session != nil {
				h.ctx.Auth.Logout(clientIO.Id)
				builder.SendPlayerLeft(clientIO.BroadcastWriter, session.Entity.Id)
			}
		}
		if h.ctx.IsClient() {
			session = h.ctx.Me
		}

		if session != nil && session.Entity != nil {
			h.ctx.Entity.Delete(session.Entity.Id) // TODO: move this elsewhere
		}
	})()

	var msg_type message.Type
	for {
		err := read(clientIO.Reader, &msg_type)
		if err == io.EOF {
			log.Println("Connection closed.")
			return
		} else if err != nil {
			//clientIO.Writer.Close()
			log.Println("binary.Read failed:", err)
			return
		}

		err = h.Handle(msg_type, clientIO)
		if err != nil {
			log.Println("Handle failed:", err)
			log.Println("Message type:", msg_type)
		}
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
