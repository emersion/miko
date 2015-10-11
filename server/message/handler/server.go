package handler

import (
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		// TODO: read exit code

		var username string
		sender := ctx.Auth.GetSession(io.Id)
		if sender != nil {
			username = sender.Username
			ctx.Auth.Logout(io.Id)
		} else {
			username = "[anonymous]"
		}

		if err := io.Writer.Close(); err != nil {
			return err
		}

		return builder.SendPlayerLeft(io.BroadcastWriter, username)
	},
	message.Types["login"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		password := readString(io.Reader)

		code := ctx.Auth.Login(io.Id, username, password)
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

		code := ctx.Auth.Register(io.Id, username, password)
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

		log.Println("Broadcasting chat message:", msg)

		var username string
		if ctx.Auth.HasSession(io.Id) {
			sender := ctx.Auth.GetSession(io.Id)
			username = sender.Username
		} else {
			username = "[anonymous]"
		}

		return builder.SendChatReceive(io.BroadcastWriter, username, msg)
	},
}
