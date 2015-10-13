package handler

import (
	"io"
	"fmt"
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

type TypeHandler func(*message.Context, *message.IO) error

type Handler struct {
	ctx *message.Context
	handlers map[message.Type]TypeHandler
}

func (h *Handler) flushEntitiesDiff(w io.Writer) error {
	pool := h.ctx.Entity.Flush()

	// TODO: broadcast only to clients who need it

	// Created entities
	for _, entity := range pool.Created {
		err := builder.SendEntityCreate(w, entity)
		if err != nil {
			return err
		}
	}

	// Updated entities
	entities := make([]*message.Entity, len(pool.Updated))
	diffs := make([]*message.EntityDiff, len(pool.Updated))
	i := 0
	for entityId, diff := range pool.Updated {
		entities[i] = h.ctx.Entity.Get(entityId)
		diffs[i] = diff
		i++
	}

	err := builder.SendEntitiesUpdate(w, entities, diffs)
	if err != nil {
		return err
	}

	// Deleted entities
	for _, entityId := range pool.Deleted {
		err := builder.SendEntityDestroy(w, entityId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Handle(t message.Type, io *message.IO) error {
	if val, ok := h.handlers[t]; ok {
		err := val(h.ctx, io)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown message type: %d", t)
	}

	if h.ctx.Type == message.ServerContext {
		// No errors, send updates
		if h.ctx.Entity.IsDirty() {
			err := h.flushEntitiesDiff(io.BroadcastWriter)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *Handler) Listen(clientIO *message.IO) {
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
