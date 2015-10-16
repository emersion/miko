package handler

import (
	"log"
	"git.emersion.fr/saucisse-royale/miko/server/message"
	"git.emersion.fr/saucisse-royale/miko/server/message/builder"
)

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["version"]: func(ctx *message.Context, io *message.IO) error {
		var version message.ProtocolVersion
		read(io.Reader, &version)

		if version != message.CURRENT_VERSION {
			exitCode := message.ExitCodes["client_outdated"]
			if version > message.CURRENT_VERSION {
				exitCode = message.ExitCodes["server_outdated"]
			}

			if err := builder.SendExit(io.Writer, exitCode); err != nil {
				return err
			}

			if err := io.Writer.Close(); err != nil {
				return err
			}
		} else {
			io.Version = version
		}

		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		// TODO: read exit code

		session := ctx.Auth.GetSession(io.Id)
		if session != nil {
			ctx.Entity.Delete(session.Entity.Id) // TODO: move this elsewhere
			ctx.Auth.Logout(io.Id)
		}

		if err := io.Writer.Close(); err != nil {
			return err
		}

		if session != nil {
			return builder.SendPlayerLeft(io.BroadcastWriter, session.Entity.Id)
		} else {
			return nil
		}
	},
	message.Types["login"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io.Reader)
		password := readString(io.Reader)

		code := ctx.Auth.Login(io.Id, username, password)
		if err := builder.SendLoginResp(io.Writer, code); err != nil {
			return err
		}

		if code == message.LoginResponseCodes["ok"] {
			log.Println("Client logged in:", username)

			session := ctx.Auth.GetSession(io.Id)
			ctx.Entity.Add(session.Entity) // TODO: move this elsewhere

			// Send initial terrain
			pos := session.Entity.Position
			err := builder.SendTerrainUpdate(io.Writer, ctx.Terrain.GetBlockAt(pos.BX, pos.BY))
			if err != nil {
				return err
			}

			return builder.SendPlayerJoined(io.BroadcastWriter, session.Entity.Id, username)
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
		var size uint8
		read(io.Reader, &size)

		for i := 0; i < int(size); i++ {
			var x, y message.BlockCoord
			read(io.Reader, &x)
			read(io.Reader, &y)

			err := builder.SendTerrainUpdate(io.Writer, ctx.Terrain.GetBlockAt(x, y))
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["action"]: func(ctx *message.Context, io *message.IO) error {
		action := &message.Action{}
		read(io.Reader, &action.Id)
		// TODO: action params

		log.Println("Client triggered action:", action.Id)

		return builder.SendActions(io.BroadcastWriter, []*message.Action{action})
	},
	message.Types["chat_send"]: func(ctx *message.Context, io *message.IO) error {
		msg := readString(io.Reader)

		log.Println("Broadcasting chat message:", msg)

		var username string
		if ctx.Auth.HasSession(io.Id) {
			session := ctx.Auth.GetSession(io.Id)
			username = session.Username
		} else {
			username = "[anonymous]"
		}

		return builder.SendChatReceive(io.BroadcastWriter, username, msg)
	},
}
