package handler

import (
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

var commonHandlers = &map[message.Type]TypeHandler{
	message.Types["ping"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Ping received!")
		return builder.SendPong(io.Writer)
	},
	message.Types["pong"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Pong received!")
		return nil
	},
}
