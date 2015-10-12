package handler

import (
	"io"
	"errors"
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
	for i, update := range pool.Updated {
		entities[i] = h.ctx.Entity.Get(update.EntityId)
		diffs[i] = update.Diff
	}

	err := builder.SendEntitiesUpdate(w, entities, diffs)
	if err != nil {
		return err
	}

	// Delete entities
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
		return errors.New("Unknown message type")
	}

	// No errors, send updates
	if h.ctx.Entity.IsDirty() {
		err := h.flushEntitiesDiff(io.BroadcastWriter)
		if err != nil {
			return err
		}
	}

	return nil
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
