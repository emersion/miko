package handler

import (
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
)

var clientHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		var code message.ExitCode
		read(io.Reader, code)
		log.Println("Server closed connection, reason:", code)

		return io.Writer.Close()
	},
	message.Types["chat_receive"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		msg := readString(io.Reader)

		log.Println("Chat:", username, msg)
		return nil
	},
}
