package handler

import (
	"errors"
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

type TypeHandler func(*message.Context, *message.IO) error

type Handler struct {
	ctx *message.Context
	handlers map[message.Type]TypeHandler
}

func (h *Handler) Handle(t message.Type, io *message.IO) error {
	if val, ok := h.handlers[t]; ok {
		return val(h.ctx, io)
	} else {
		return errors.New("Unknown message type")
	}
}

func (h *Handler) Listen(io *message.IO) {
	var msg_type message.Type
	for {
		err := read(io.Reader, &msg_type)
		if err != nil {
			io.Writer.Close()
			log.Println("binary.Read failed:", err)
			return
		}

		err = h.Handle(msg_type, io)
		if err != nil {
			log.Println("Handle failed:", err)
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
	if ctx.Type == message.ServerContext {
		handlers = mergeHandlers(commonHandlers, serverHandlers)
	} else if ctx.Type == message.ClientContext {
		handlers = mergeHandlers(commonHandlers, clientHandlers)
	} else {
		handlers = *commonHandlers
	}

	return &Handler{
		ctx: ctx,
		handlers: handlers,
	}
}
