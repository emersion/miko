package handler

import (
	"io"
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

func ReadExit(r io.Reader) (code message.ExitCode) {
	read(r, &code)
	return
}

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
