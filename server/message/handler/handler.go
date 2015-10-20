package handler

import (
	"io"
	"fmt"
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

// A handler for a specific message
type TypeHandler func(*message.Context, *message.IO) error

// Handles messages from remote
type Handler struct {
	ctx *message.Context
	handlers map[message.Type]TypeHandler
}

// Handle a message of the specified type
func (h *Handler) Handle(t message.Type, io *message.IO) error {
	if h.ctx.IsServer() {
		// TODO: Check that the client sent his version
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
			err := builder.SendEntitiesDiffToClients(io.BroadcastWriter, h.ctx.Entity.Flush())
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

		session := h.ctx.Auth.GetSession(clientIO.Id)
		if session != nil {
			h.ctx.Entity.Delete(session.Entity.Id) // TODO: move this elsewhere

			if h.ctx.IsServer() {
				h.ctx.Auth.Logout(clientIO.Id)
				builder.SendPlayerLeft(clientIO.BroadcastWriter, session.Entity.Id)
			}
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
		ctx: ctx,
		handlers: handlers,
	}
}
