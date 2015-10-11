package handler

import (
	"errors"
	"log"
	".."
)

type TypeHandler func(*message.Context, *message.IO) error

type Handler struct {
	ctx *message.Context
	handlers map[message.Type]TypeHandler
}

func (h *Handler) Handle(t message.Type, io *message.IO) error {
	if val, ok := h.handlers[t]; ok {
		return val(h.ctx, io);
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
		h.Handle(msg_type, io)
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
	return &Handler{
		ctx: ctx,
		handlers: mergeHandlers(commonHandlers, serverHandlers),
	}
}
