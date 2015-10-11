package handler

import (
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		msg := readString(io.Reader)

		log.Println("Chat:", username, msg)
		return nil
	},
}
