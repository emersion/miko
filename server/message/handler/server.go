package handler

import (
	"log"
	".."
	"../builder"
)

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["ping"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Ping received!")
		return builder.SendPong(io.Writer)
	},
	message.Types["pong"]: func(ctx *message.Context, io *message.IO) error {
		log.Println("Pong received!")
		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		sender := ctx.Auth.GetSession(io.Id)

		// TODO: check if user is logged in
		ctx.Auth.Logout(io)

		if err := io.Writer.Close(); err != nil {
			return err
		}

		return builder.SendPlayerLeft(io.BroadcastWriter, sender.Username)
	},
	message.Types["login"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		password := readString(io.Reader)

		code := ctx.Auth.Login(io, username, password)
		if err := builder.SendLoginResp(io.Writer, code); err != nil {
			return err
		}

		if code == message.LoginResponseCodes["ok"] {
			return builder.SendPlayerJoined(io.BroadcastWriter, username)
		} else {
			return nil
		}
	},
	message.Types["register"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		password := readString(io.Reader)

		code := ctx.Auth.Register(io, username, password)
		return builder.SendRegisterResp(io.Writer, code)
	},
	message.Types["terrain_request"]: func(ctx *message.Context, io *message.IO) error {
		var x, y message.BlockCoord
		read(io.Reader, x)
		read(io.Reader, y)

		return builder.SendTerrainUpdate(io.Writer, ctx.Terrain.GetBlockAt(x, y))
	},
	message.Types["chat_send"]: func(ctx *message.Context, io *message.IO) error {
		msg := readString(io.Reader)
		sender := ctx.Auth.GetSession(io.Id)

		return builder.SendChatReceive(io.BroadcastWriter, sender.Username, msg)
	},
}
