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
		read(io, &version)
		io.Version = version

		if version != message.CurrentVersion {
			code := message.ExitCodes["client_outdated"]
			if version > message.CurrentVersion {
				code = message.ExitCodes["server_outdated"]
			}

			if err := builder.SendExit(io, code); err != nil {
				return err
			}

			if err := io.Close(); err != nil {
				return err
			}
		} else {
			if err := builder.SendConfig(io, ctx.Config); err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, io *message.IO) error {
		code := ReadExit(io)
		log.Println("Client disconnected with code:", code)

		if err := io.Close(); err != nil {
			return err
		}

		return nil
	},
	message.Types["login"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io)
		password := readString(io)

		code := ctx.Auth.Login(io.Id, username, password)
		if err := builder.SendLoginResp(io, code, ctx.Clock.GetRelativeTick()); err != nil {
			return err
		}

		log.Println("Client logged in:", username, code)

		if code == message.LoginResponseCodes["ok"] {
			// Create entity
			session := ctx.Auth.GetSession(io.Id)
			if session == nil {
				return errors.New("Cannot get newly logged in user's session")
			}

			session.Entity.Type = 0   // player
			session.Entity.Sprite = 1 // player

			ctx.Entity.Add(session.Entity, ctx.Clock.GetAbsoluteTick()) // TODO: move this elsewhere

			// Send initial terrain

			pos := session.Entity.Position
			blk, err := ctx.Terrain.GetBlockAt(pos.BX, pos.BY)
			if err != nil {
				return err
			}

			err = builder.SendTerrainUpdate(io, ctx.Clock.GetRelativeTick(), blk)
			if err != nil {
				return err
			}

			// Broadcast new entity
			err = builder.SendEntitiesDiffToClients(io.Broadcaster(), ctx.Clock.GetRelativeTick(), ctx.Entity.Flush())
			if err != nil {
				return err
			}

			return builder.SendPlayerJoined(io.Broadcaster(), ctx.Clock.GetRelativeTick(), session.Entity.Id, username)
		} else {
			return nil
		}
	},
	message.Types["register"]: func(ctx *message.Context, io *message.IO) error {
		username := readString(io)
		password := readString(io)

		code := ctx.Auth.Register(io.Id, username, password)

		log.Println("Client registered:", username, code)

		return builder.SendRegisterResp(io, code)
	},
	message.Types["terrain_request"]: func(ctx *message.Context, io *message.IO) error {
		var size uint8
		read(io, &size)

		for i := 0; i < int(size); i++ {
			var x, y message.BlockCoord
			read(io, &x)
			read(io, &y)

			blk, err := ctx.Terrain.GetBlockAt(x, y)
			if err != nil {
				return err
			}

			err = builder.SendTerrainUpdate(io, ctx.Clock.GetRelativeTick(), blk)
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["entity_update"]: func(ctx *message.Context, io *message.IO) error {
		t := ctx.Clock.ToAbsoluteTick(readTick(io))

		entity, diff := ReadEntity(io)
		ctx.Entity.Update(entity, diff, t)

		return nil
	},
	message.Types["action_do"]: func(ctx *message.Context, io *message.IO) error {
		if !ctx.Auth.HasSession(io.Id) {
			// TODO: trigger an error
			return nil
		}

		session := ctx.Auth.GetSession(io.Id)

		action := &message.Action{
			Initiator: session.Entity.Id,
		}

		readTick(io) // TODO: properly handle this tick
		read(io, &action.Id)
		// TODO: action params

		log.Println("Client triggered action:", action.Id)

		return builder.SendActionsDone(io.Broadcaster(), ctx.Clock.GetRelativeTick(), []*message.Action{action})
	},
	message.Types["chat_send"]: func(ctx *message.Context, io *message.IO) error {
		msg := readString(io)

		if !ctx.Auth.HasSession(io.Id) {
			return errors.New("User not authenticated")
		}

		session := ctx.Auth.GetSession(io.Id)
		username := session.Username

		log.Println("Broadcasting chat message:", username, msg)

		return builder.SendChatReceive(io.Broadcaster(), ctx.Clock.GetRelativeTick(), username, msg)
	},
}
