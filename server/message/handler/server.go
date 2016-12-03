package handler

import (
	"errors"
	"git.emersion.fr/saucisse-royale/miko.git/server/message"
	"git.emersion.fr/saucisse-royale/miko.git/server/message/builder"
	"io"
	"log"
	"time"
)

func ReadVersion(r io.Reader) (version message.ProtocolVersion) {
	Read(r, &version)
	return
}

func ReadLogin(r io.Reader) (username, password string) {
	Read(r, &username, &password)
	return
}

func ReadRegister(r io.Reader) (username, password string) {
	Read(r, &username, &password)
	return
}

func ReadChatSend(r io.Reader) (msg string) {
	Read(r, &msg)
	return
}

var serverHandlers = &map[message.Type]TypeHandler{
	message.Types["version"]: func(ctx *message.Context, conn *message.Conn) error {
		conn.Version = ReadVersion(conn)

		if conn.Version != message.CurrentVersion {
			code := message.ExitCodes["client_outdated"]
			if conn.Version > message.CurrentVersion {
				code = message.ExitCodes["server_outdated"]
			}

			if err := builder.SendExit(conn, code); err != nil {
				return err
			}

			if err := conn.Close(); err != nil {
				return err
			}
		} else {
			conn.State = message.Accepted

			if err := builder.SendConfig(conn, ctx.Config); err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["exit"]: func(ctx *message.Context, conn *message.Conn) error {
		code := ReadExit(conn)
		log.Println("Client disconnected with code:", code)

		if err := conn.Close(); err != nil {
			return err
		}

		return nil
	},
	message.Types["login"]: func(ctx *message.Context, conn *message.Conn) error {
		username, password := ReadLogin(conn)

		code := ctx.Auth.Login(conn.Id, username, password)
		if err := builder.SendLoginResp(conn, code, ctx.Clock.GetRelativeTick(), time.Now()); err != nil {
			return err
		}

		log.Println("Client requested login:", username, code)

		if code != message.LoginResponseCodes["ok"] {
			return nil
		}

		conn.State = message.LoggedIn

		// Get entity
		session := ctx.Auth.GetSession(conn.Id)
		if session == nil {
			return errors.New("Cannot get newly logged in user's session")
		}

		// Send initial terrain
		log.Println("Sending initial terrain")
		pos := session.Entity.Position
		radius := message.BlockCoord(8)
		blks := []*message.Block{}
		for i := pos.BX - radius; i <= pos.BX+radius; i++ {
			for j := pos.BY - radius; j <= pos.BY+radius; j++ {
				blk, err := ctx.Terrain.GetBlockAt(i, j)
				if err != nil {
					log.Println("Cannot send initial terrain to client", err)
					continue
				}

				blks = append(blks, blk)
			}
		}
		err := builder.SendChunksUpdate(conn, ctx.Clock.GetRelativeTick(), blks)
		if err != nil {
			return err
		}

		// Send initial entities
		log.Println("Sending initial entities")
		for _, e := range ctx.Entity.List() {
			if e.Id == session.Entity.Id {
				continue
			}

			err := builder.SendEntityCreate(conn, ctx.Clock.GetRelativeTick(), e)
			if err != nil {
				return err
			}
		}

		// Send current players
		for _, s := range ctx.Auth.List() {
			if s.Id == session.Id {
				continue
			}

			err := builder.SendPlayerJoined(conn, ctx.Clock.GetRelativeTick(), s.Entity.Id, s.Username)
			if err != nil {
				return err
			}
		}

		// Register new entity
		req := ctx.Entity.Add(session.Entity, ctx.Clock.GetAbsoluteTick()) // TODO: move this elsewhere?
		err = req.Wait()
		log.Println("Entity registered!")
		if err != nil {
			return err
		}

		// Broadcast new entity to other clients
		log.Println("Flushing entities diff")
		err = builder.SendEntitiesDiffToClients(conn.Broadcaster(), ctx.Clock.GetRelativeTick(), ctx.Entity.Flush())
		if err != nil {
			return err
		}

		// Send new entity to this client
		err = builder.SendEntityCreate(conn, ctx.Clock.GetRelativeTick(), session.Entity)
		if err != nil {
			return err
		}

		// Mark the io as ready
		// It will now receive all broadcasts
		conn.State = message.Ready

		// Broadcast player_joined to everyone (including this client)
		err = builder.SendPlayerJoined(conn.Broadcaster(), ctx.Clock.GetRelativeTick(), session.Entity.Id, username)
		if err != nil {
			return err
		}

		return nil
	},
	message.Types["register"]: func(ctx *message.Context, conn *message.Conn) error {
		username, password := ReadRegister(conn)

		code := ctx.Auth.Register(conn.Id, username, password)

		log.Println("Client registered:", username, code)

		return builder.SendRegisterResp(conn, code)
	},
	message.Types["terrain_request"]: func(ctx *message.Context, conn *message.Conn) error {
		var size uint8
		Read(conn, &size)

		for i := 0; i < int(size); i++ {
			var x, y message.BlockCoord
			Read(conn, &x, &y)

			blk, err := ctx.Terrain.GetBlockAt(x, y)
			if err != nil {
				return err
			}

			err = builder.SendChunkUpdate(conn, ctx.Clock.GetRelativeTick(), blk)
			if err != nil {
				return err
			}
		}

		return nil
	},
	message.Types["entity_update"]: func(ctx *message.Context, conn *message.Conn) error {
		// TODO: add update initiator as parameter to ctx.Entity.Update()

		t := ctx.Clock.ToAbsoluteTick(readTick(conn))

		entity, diff := ReadEntity(conn)
		ctx.Entity.Update(entity, diff, t)

		return nil
	},
	message.Types["action_do"]: func(ctx *message.Context, conn *message.Conn) error {
		if !ctx.Auth.HasSession(conn.Id) {
			return errors.New("Cannot execute action: remote not logged in")
		}

		session := ctx.Auth.GetSession(conn.Id)

		action := &message.Action{
			Initiator: session.Entity.Id,
		}

		t := ctx.Clock.ToAbsoluteTick(readTick(conn))
		Read(conn, &action.Id)

		// TODO: move action params somewhere else
		// [GAME-SPECIFIC]
		if action.Id == 0 { // throw_ball
			var angle float32
			var tmpId message.EntityId
			Read(conn, &angle, &tmpId)
			log.Println("Received ball action:", angle, tmpId)

			action.Params = []interface{}{angle, tmpId}
		} else {
			log.Println("Client triggered unknown action:", action.Id)
		}

		ctx.Action.Execute(action, t)
		return nil
	},
	message.Types["chat_send"]: func(ctx *message.Context, conn *message.Conn) error {
		msg := ReadChatSend(conn)

		if !ctx.Auth.HasSession(conn.Id) {
			return errors.New("User not authenticated")
		}

		session := ctx.Auth.GetSession(conn.Id)
		username := session.Username

		log.Println("Broadcasting chat message:", username, msg)

		return builder.SendChatReceive(conn.Broadcaster(), ctx.Clock.GetRelativeTick(), username, msg)
	},
}
