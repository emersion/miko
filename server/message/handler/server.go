package handler

import (
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"log"
)

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["version"]: func(ctx *message.Context, io *message.IO) error {
		var version message.ProtocolVersion
		read(io.Reader, &version)
		io.Version = version

		code := message.VersionResponseCodes["ok"]

		if version != message.CurrentVersion {
			code = message.VersionResponseCodes["client_outdated"]
			if version > message.CurrentVersion {
				code = message.VersionResponseCodes["server_outdated"]
			}
		}

		if err := builder.SendVersionResponse(io.Writer, code); err != nil {
			return err
		}

		if code != message.VersionResponseCodes["ok"] {
			if err := io.Writer.Close(); err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadExit(io.Reader)
		log.Println("Client disconnected with code:", code)

		if err := io.Writer.Close(); err != nil {
			return err
		}

		return nil
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

			// Create entity
			session := ctx.Auth.GetSession(io.Id)
			if session == nil {
				return errors.New("Cannot get newly logged in user's session")
			}
			ctx.Entity.Add(session.Entity) // TODO: move this elsewhere

			// Send initial terrain

			pos := session.Entity.Position
			blk, err := ctx.Terrain.GetBlockAt(pos.BX, pos.BY)
			if err != nil {
				return err
			}

			err = builder.SendTerrainUpdate(io.Writer, ctx.Clock.GetTick(), blk)
			if err != nil {
				return err
			}

			// Broadcast new entity
			err = builder.SendEntitiesDiffToClients(io.BroadcastWriter, ctx.Clock.GetTick(), ctx.Entity.Flush())
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

			blk, err := ctx.Terrain.GetBlockAt(x, y)
			if err != nil {
				return err
			}

			err = builder.SendTerrainUpdate(io.Writer, ctx.Clock.GetTick(), blk)
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["entity_update"]: func(ctx *message.Context, io *message.IO) error {
		// TODO: security checks

		readTick(io.Reader)

		entity, diff := ReadEntity(io.Reader)
		ctx.Entity.Update(entity, diff)

		return nil
	},
	message.Types["action_do"]: func(ctx *message.Context, io *message.IO) error {
		if !ctx.Auth.HasSession(io.Id) {
			// TODO: trigger an error
			return nil
		}

		readTick(io.Reader)

		session := ctx.Auth.GetSession(io.Id)

		action := &message.Action{
			Initiator: session.Entity.Id,
		}
		read(io.Reader, &action.Id)
		// TODO: action params

		log.Println("Client triggered action:", action.Id)

		return builder.SendActionsDone(io.BroadcastWriter, ctx.Clock.GetTick(), []*message.Action{action})
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
